package store

import (
	"path/filepath"
	"testing"
	"time"
)

func open(t *testing.T) *Store {
	t.Helper()
	s, err := Open(filepath.Join(t.TempDir(), "t.db"))
	if err != nil {
		t.Fatal(err)
	}
	return s
}

// A mail header (no body yet) has NULL body columns; listing must not error.
func TestListMailsHandlesHeaderOnly(t *testing.T) {
	s := open(t)
	defer s.Close()
	if _, err := s.InsertMailInfo(Mail{ID: 1, AlbionServerID: 1, AuctionType: 1, RawLocationID: "3005", PlayerName: "Hero", Received: time.Now(), Type: 1}); err != nil {
		t.Fatal(err)
	}
	mails, err := s.ListMails(100)
	if err != nil {
		t.Fatalf("ListMails: %v", err)
	}
	if len(mails) != 1 || mails[0].ID != 1 || mails[0].PartialAmount != 0 {
		t.Fatalf("unexpected: %+v", mails)
	}
	// now read it (sets body) and confirm values come through
	if err := s.SetMailData(1, 5, 10, "T4_BAG", 154984, 0, 30996.8); err != nil {
		t.Fatal(err)
	}
	mails, _ = s.ListMails(100)
	if mails[0].ItemID != "T4_BAG" || mails[0].PartialAmount != 5 || mails[0].UnitSilver != 30996.8 {
		t.Fatalf("after SetMailData: %+v", mails[0])
	}
}

func TestTradesRoundTrip(t *testing.T) {
	s := open(t)
	defer s.Close()
	if err := s.InsertTrade(Trade{ID: "g1", AlbionServerID: 1, Amount: 2, DateTime: time.Now(), ItemID: "T4_BAG", UnitSilver: 100}); err != nil {
		t.Fatal(err)
	}
	tr, err := s.ListTrades(10)
	if err != nil || len(tr) != 1 || tr[0].ID != "g1" {
		t.Fatalf("trades: %v %+v", err, tr)
	}
}
