# Albion Market Data Client

A Windows desktop app for **[Albion Online](https://albiononline.com)** that passively reads the
game's own network traffic to collect **market prices** and **gameplay statistics**, and syncs them to
your account at **[albionmarket.gg](https://albionmarket.gg)**.

It does **not** modify the game, inject packets, or automate anything. It only *reads* the game's
Photon UDP packets (the same data your client already sends and receives) and uploads the parts you
opt into.

---

## What it does

- **Market prices** — captures marketplace buy/sell orders, price history, item value estimates and
  gold prices as you browse the in-game market, and contributes them to the shared price data.
- **Gameplay tracking** (optional, per-feature toggles):
  - **Loot** — items you and your party pick up.
  - **Damage meter** — damage dealt/taken, per-encounter and group session totals, with a fame chart.
  - **Gathering** — resources gathered and their value.
  - **Awakened weapons** — awakened rolls, synced to your dashboard.
  - **Destiny Board / specs**, **trades**, **mails**, **party** composition.
- **Live dashboard** — sessions, capture feed, and stats in a native window (system tray + optional
  start-with-Windows).

---

## How capture works (and what it does *not* see)

Capture uses **[Npcap](https://npcap.com)**, a packet-capture driver. Npcap is technically capable of
reading all network traffic — **this app does not.** It applies a filter and only reads the game's
**Photon UDP** packets. It does **not** read, store, or transmit your browsing, other applications, or
any non-game network traffic.

It also collects **no** passwords, keystrokes, screen contents, clipboard, microphone, or files. Login
uses Discord's OAuth device flow, so the app never sees your Discord password.

You don't have to take our word for it — the client source is in this repository.

---

## Privacy & data

- Gameplay uploads happen **only when you're signed in**, and each category is an **independent toggle
  you can turn off**.
- If you don't sign in, no gameplay or account data is uploaded.
- The Discord refresh token stored on your PC is **encrypted at rest** (Windows DPAPI, tied to your
  Windows account).
- Note: combat/loot/party data can include other players' **in-game character names** (party members,
  nearby looters) — this is public in-game information.

A full **privacy & data-collection policy** is published on the website:
**[albionmarket.gg](https://albionmarket.gg)**.

## Security

- All traffic to the backend is **HTTPS with full certificate validation** — no disabled certificate
  checks, no plain-HTTP fallback.
- The updater is **notify-only**: it checks a version endpoint and, if an update exists, opens a
  download link in your browser. It **never** silently downloads or runs anything.
- The installer is currently **unsigned**, so Windows SmartScreen / some browsers may warn on download
  (code signing is being added). See below.

---

## Install

1. Download the latest `AlbionMarketDataClient-Setup.exe` from the
   [Releases](https://github.com/niick1231/albionmarket_dataclient/releases) page.
2. Run it. The installer will offer to install **Npcap** if it isn't already present (required for
   packet capture).
3. Launch the app and sign in with Discord to enable gameplay syncing.

> **Unsigned-binary warning:** because the installer isn't code-signed yet, SmartScreen may show a
> "Windows protected your PC" prompt and some browsers may warn on download. This is expected for a new,
> unsigned binary; choose *More info → Run anyway* if you trust the source. Signing is on the roadmap.

---

## Build from source

Requirements: **Go 1.26+**, **[Wails v2](https://wails.io)**, a C toolchain (cgo), and the **Npcap
SDK** (capture uses gopacket + Npcap).

```sh
# from the repo root
source scripts/env.sh        # sets up the cgo / Npcap SDK environment
wails build                  # dev build → build/bin/
# release installer:
wails build -platform windows/amd64 -nsis -ldflags "-X main.version=<version>"
```

The module path is `github.com/niick1231/albionmarket_dataclient`. The frontend lives in `frontend/`
(TypeScript + Vite), the Go pipeline in `internal/`.

---

## Configuration

By default the client points at `https://albionmarket.gg` for both ingest and auth. Settings (capture
device, upload toggles, tray/start-up behaviour, etc.) are configurable in the app. Endpoints can be
overridden in the local config file for self-hosting.

---

## Credits & acknowledgements

- The **Photon protocol parsing**, event/operation code tables, and the proof-of-work upload handshake
  derive from the wider **[Albion Online Data Project](https://www.albion-online-data.com)** ecosystem
  and the **AlbionDataAvalonia** project. Thanks to those projects and their contributors.
- **[Npcap](https://npcap.com)** (the Nmap Project) — packet capture.
- **[Wails](https://wails.io)** — the Go + web desktop framework.
- **[gopacket](https://github.com/gopacket/gopacket)** — packet decoding.
- Game reference data from the community **`ao-bin-dumps`**.

Albion Online is a trademark of Sandbox Interactive GmbH. This project is a fan-made tool and is not
affiliated with or endorsed by Sandbox Interactive.

---

## License

License to be finalized before public release — see `LICENSE`.
