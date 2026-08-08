// Package capture sniffs Albion's Photon UDP traffic and feeds it to the parser.
package capture

import (
	"strings"
	"sync"
	"time"

	"github.com/gopacket/gopacket"
	"github.com/gopacket/gopacket/layers"
	"github.com/gopacket/gopacket/pcap"

	"github.com/albionmarketgg/dataclient/internal/config"
	"github.com/albionmarketgg/dataclient/internal/photon"
	"github.com/albionmarketgg/dataclient/internal/state"
)

// rescanInterval is how often the watchdog re-enumerates capture devices while
// running. Adapters routinely appear AFTER capture starts (autostart racing
// Wi-Fi/DHCP at boot, docking, VPN up/down); without rescanning, a client
// started before the network is up captures nothing until manually restarted.
const rescanInterval = 20 * time.Second

// Listener captures packets across all devices and dispatches Photon payloads.
type Listener struct {
	cfg    config.Config
	st     *state.State
	parser *photon.Parser
	logf   func(string)

	mu       sync.Mutex
	handles  map[string]*pcap.Handle // keyed by pcap device name
	running  bool
	stopOnce chan struct{}

	narrowed   bool
	activeName string
}

// New builds a Listener.
func New(cfg config.Config, st *state.State, parser *photon.Parser, logf func(string)) *Listener {
	if logf == nil {
		logf = func(string) {}
	}
	return &Listener{cfg: cfg, st: st, parser: parser, logf: logf}
}

// Devices returns the available capture device descriptions.
func Devices() ([]string, error) {
	devs, err := pcap.FindAllDevs()
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(devs))
	for _, d := range devs {
		name := d.Description
		if name == "" {
			name = d.Name
		}
		out = append(out, name)
	}
	return out, nil
}

// Running reports capture status.
func (l *Listener) Running() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.running
}

// Start opens all devices, begins capturing, and keeps a watchdog rescanning
// for devices that appear later (e.g. the client autostarted before the network
// was up). Non-blocking; only a missing pcap runtime is a hard error.
func (l *Listener) Start() error {
	l.mu.Lock()
	if l.running {
		l.mu.Unlock()
		return nil
	}
	l.running = true
	l.narrowed = false
	l.handles = map[string]*pcap.Handle{}
	l.stopOnce = make(chan struct{})
	stop := l.stopOnce
	l.mu.Unlock()

	if d := l.cfg.NetworkStartDelaySecs; d > 0 {
		time.Sleep(time.Duration(d) * time.Second)
	}

	opened, err := l.scan(stop)
	if err != nil {
		l.setStopped()
		return err
	}

	if opened == 0 {
		l.logf("No capture devices available yet - will keep retrying (is the network up and a capture driver installed?).")
	} else {
		l.logf("Listening for Albion market traffic on " + itoa(opened) + " device(s).")
		l.st.SetListening(true)
	}
	go l.watchdog(stop)
	return nil
}

// scan opens every eligible capture device that isn't open yet. Returns how
// many new devices were opened; the only error is device enumeration failing
// outright (missing pcap runtime).
func (l *Listener) scan(stop chan struct{}) (int, error) {
	devs, err := pcap.FindAllDevs()
	if err != nil {
		return 0, err
	}
	opened := 0
	for _, dev := range devs {
		if want := strings.TrimSpace(l.cfg.CaptureDevice); want != "" {
			if !strings.Contains(dev.Description, want) && !strings.Contains(dev.Name, want) {
				continue
			}
		}
		l.mu.Lock()
		_, already := l.handles[dev.Name]
		l.mu.Unlock()
		if already {
			continue
		}
		h, err := pcap.OpenLive(dev.Name, 65536, false, pcap.BlockForever)
		if err != nil {
			continue
		}
		if err := h.SetBPFFilter(l.cfg.PacketFilter); err != nil {
			h.Close()
			continue
		}
		l.mu.Lock()
		l.handles[dev.Name] = h
		l.mu.Unlock()
		opened++
		go l.readLoop(dev.Name, h, dev.Description, stop)
	}
	return opened, nil
}

// watchdog periodically rescans so adapters that appear after startup (late
// Wi-Fi/DHCP at boot, docking, VPN) get captured without an app restart.
func (l *Listener) watchdog(stop chan struct{}) {
	t := time.NewTicker(rescanInterval)
	defer t.Stop()
	for {
		select {
		case <-stop:
			return
		case <-t.C:
			l.mu.Lock()
			had := len(l.handles)
			l.mu.Unlock()
			opened, err := l.scan(stop)
			if err != nil || opened == 0 {
				continue
			}
			if had == 0 {
				l.logf("Capture device available - listening for Albion market traffic on " + itoa(opened) + " device(s).")
				l.st.SetListening(true)
			} else {
				l.logf("Capture attached to " + itoa(opened) + " new device(s).")
			}
		}
	}
}

func (l *Listener) readLoop(name string, h *pcap.Handle, desc string, stop chan struct{}) {
	src := gopacket.NewPacketSource(h, h.LinkType())
	src.NoCopy = true
	in := src.Packets()
	for {
		select {
		case <-stop:
			return
		case pkt, ok := <-in:
			if !ok {
				// The device died (disabled/unplugged/driver reset). Forget it so
				// the watchdog can reattach if it comes back.
				l.mu.Lock()
				if l.handles != nil && l.handles[name] == h {
					delete(l.handles, name)
				}
				l.mu.Unlock()
				h.Close()
				return
			}
			l.handlePacket(pkt, desc)
		}
	}
}

func (l *Listener) handlePacket(pkt gopacket.Packet, desc string) {
	udpLayer := pkt.Layer(layers.LayerTypeUDP)
	if udpLayer == nil {
		return
	}
	udp, _ := udpLayer.(*layers.UDP)
	if udp == nil || len(udp.Payload) == 0 {
		return
	}
	l.st.MarkPacket()

	// narrow to the device that delivered the first valid packet
	l.mu.Lock()
	if !l.narrowed {
		l.narrowed = true
		l.activeName = desc
		go l.closeOthers(desc)
	}
	l.mu.Unlock()

	if ipLayer := pkt.Layer(layers.LayerTypeIPv4); ipLayer != nil {
		ip, _ := ipLayer.(*layers.IPv4)
		if ip != nil {
			l.detectServer(ip.SrcIP.String())
		}
	}

	status := l.parser.ReceivePacket(udp.Payload)
	switch status {
	case photon.StatusEncrypted:
		l.st.SetEncrypted(true)
	case photon.StatusSuccess:
		l.st.SetEncrypted(false)
	}
}

func (l *Listener) detectServer(srcIP string) {
	for i := range l.cfg.Servers {
		srv := &l.cfg.Servers[i]
		for _, prefix := range srv.HostIPs {
			if strings.HasPrefix(srcIP, prefix) {
				l.st.SetServer(srv)
				return
			}
		}
	}
}

func (l *Listener) closeOthers(keep string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	// Best-effort: per-handle desc isn't tracked, so keep all handles open.
	_ = keep
}

// Stop ends capture and closes all handles.
func (l *Listener) Stop() {
	l.mu.Lock()
	if !l.running {
		l.mu.Unlock()
		return
	}
	l.running = false
	close(l.stopOnce)
	handles := l.handles
	l.handles = nil
	l.mu.Unlock()

	for _, h := range handles {
		h.Close()
	}
	l.st.SetListening(false)
	l.logf("Stopped listening.")
}

func (l *Listener) setStopped() {
	l.mu.Lock()
	l.running = false
	l.mu.Unlock()
	l.st.SetListening(false)
}
