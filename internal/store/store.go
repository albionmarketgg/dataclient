// Package store provides local SQLite persistence for mails and trades.
package store

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

// Store wraps the SQLite database.
type Store struct {
	db *sql.DB
}

// Open opens (creating if needed) the database at path and ensures the schema.
func Open(path string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", path+"?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)")
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}

// Close closes the database.
func (s *Store) Close() error { return s.db.Close() }

func (s *Store) migrate() error {
	const schema = `
CREATE TABLE IF NOT EXISTS user_auth (
  UserId TEXT PRIMARY KEY,
  RefreshToken TEXT NOT NULL,
  Username TEXT,
  Avatar TEXT
);

CREATE TABLE IF NOT EXISTS mails (
  Id INTEGER PRIMARY KEY,
  AlbionServerId INTEGER,
  AuctionType INTEGER,
  Deleted INTEGER DEFAULT 0,
  IsSet INTEGER DEFAULT 0,
  ItemId TEXT,
  LocationId INTEGER,
  RawLocationId TEXT,
  PlayerName TEXT,
  PartialAmount INTEGER,
  TotalAmount INTEGER,
  Received TEXT,
  TaxesPercent REAL,
  TotalSilver INTEGER,
  TotalTaxes INTEGER,
  Type INTEGER,
  UnitSilver REAL,
  Synced INTEGER DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_mails_received ON mails(AlbionServerId, LocationId, AuctionType, Deleted, Received);

CREATE TABLE IF NOT EXISTS trades (
  Id TEXT PRIMARY KEY,
  AlbionServerId INTEGER,
  Amount INTEGER,
  DateTime TEXT,
  Deleted INTEGER DEFAULT 0,
  ItemId TEXT,
  LocationId INTEGER,
  Operation INTEGER,
  PlayerName TEXT,
  QualityLevel INTEGER,
  RawLocationId TEXT,
  SalesTaxesPercent REAL,
  Type INTEGER,
  UnitSilver REAL,
  Synced INTEGER DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_trades_dt ON trades(AlbionServerId, LocationId, Deleted, DateTime);
`
	if _, err := s.db.Exec(schema); err != nil {
		return err
	}
	// Backfill the Synced column on older DBs (ignore "duplicate column").
	s.db.Exec(`ALTER TABLE trades ADD COLUMN Synced INTEGER DEFAULT 0`)
	s.db.Exec(`ALTER TABLE mails ADD COLUMN Synced INTEGER DEFAULT 0`)
	return nil
}

// ---- User auth ----

// UserAuth is the persisted login (refresh token encrypted by the auth layer).
type UserAuth struct {
	UserID       string
	RefreshToken string
	Username     string
	Avatar       string
}

// SaveUserAuth upserts the stored login (single row).
func (s *Store) SaveUserAuth(u UserAuth) error {
	_, err := s.db.Exec(`INSERT OR REPLACE INTO user_auth (UserId, RefreshToken, Username, Avatar) VALUES (?,?,?,?)`,
		u.UserID, u.RefreshToken, u.Username, u.Avatar)
	return err
}

// LoadUserAuth returns the stored login, if any.
func (s *Store) LoadUserAuth() (UserAuth, bool, error) {
	var u UserAuth
	err := s.db.QueryRow(`SELECT UserId, RefreshToken, Username, Avatar FROM user_auth LIMIT 1`).
		Scan(&u.UserID, &u.RefreshToken, &u.Username, &u.Avatar)
	if err == sql.ErrNoRows {
		return UserAuth{}, false, nil
	}
	return u, err == nil, err
}

// ClearUserAuth deletes the stored login.
func (s *Store) ClearUserAuth() error {
	_, err := s.db.Exec(`DELETE FROM user_auth`)
	return err
}

// ---- Mails ----

// Mail is a persisted marketplace summary mail.
type Mail struct {
	ID             int64     `json:"id"`
	AlbionServerID int       `json:"albionServerId"`
	AuctionType    int       `json:"auctionType"`
	Deleted        bool      `json:"deleted"`
	IsSet          bool      `json:"isSet"`
	ItemID         string    `json:"itemId"`
	LocationID     int       `json:"locationId"`
	RawLocationID  string    `json:"rawLocationId"`
	PlayerName     string    `json:"playerName"`
	PartialAmount  int       `json:"partialAmount"`
	TotalAmount    int       `json:"totalAmount"`
	Received       time.Time `json:"received"`
	TaxesPercent   float64   `json:"taxesPercent"`
	TotalSilver    int64     `json:"totalSilver"`
	TotalTaxes     int64     `json:"totalTaxes"`
	Type           int       `json:"type"`
	UnitSilver     float64   `json:"unitSilver"`
}

// InsertMailInfo inserts a mail header if its Id is new. Returns true if inserted.
func (s *Store) InsertMailInfo(m Mail) (bool, error) {
	res, err := s.db.Exec(`INSERT OR IGNORE INTO mails
		(Id, AlbionServerId, AuctionType, Deleted, IsSet, RawLocationId, LocationId, PlayerName, Received, Type, TaxesPercent)
		VALUES (?,?,?,0,0,?,?,?,?,?,0)`,
		m.ID, m.AlbionServerID, m.AuctionType, m.RawLocationID, m.LocationID, m.PlayerName, m.Received.UTC().Format(time.RFC3339Nano), m.Type)
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n > 0, nil
}

// MailIsSet reports whether the mail body has already been parsed.
func (s *Store) MailIsSet(id int64) (exists, isSet bool, err error) {
	var v int
	err = s.db.QueryRow(`SELECT IsSet FROM mails WHERE Id=?`, id).Scan(&v)
	if err == sql.ErrNoRows {
		return false, false, nil
	}
	if err != nil {
		return false, false, err
	}
	return true, v != 0, nil
}

// SetMailData fills the parsed body fields and marks the mail set.
func (s *Store) SetMailData(id int64, partial, total int, itemID string, totalSilver, totalTaxes int64, unitSilver float64) error {
	_, err := s.db.Exec(`UPDATE mails SET IsSet=1, PartialAmount=?, TotalAmount=?, ItemId=?, TotalSilver=?, TotalTaxes=?, UnitSilver=? WHERE Id=?`,
		partial, total, itemID, totalSilver, totalTaxes, unitSilver, id)
	return err
}

// GetMail loads a mail by id.
func (s *Store) GetMail(id int64) (Mail, bool, error) {
	row := s.db.QueryRow(`SELECT Id,AlbionServerId,AuctionType,Deleted,IsSet,ItemId,LocationId,RawLocationId,PlayerName,PartialAmount,TotalAmount,Received,TaxesPercent,TotalSilver,TotalTaxes,Type,UnitSilver FROM mails WHERE Id=?`, id)
	m, err := scanMail(row)
	if err == sql.ErrNoRows {
		return Mail{}, false, nil
	}
	return m, err == nil, err
}

// ListMails returns up to limit most-recent non-deleted mails.
func (s *Store) ListMails(limit int) ([]Mail, error) {
	if limit <= 0 {
		limit = 1000
	}
	rows, err := s.db.Query(`SELECT Id,AlbionServerId,AuctionType,Deleted,IsSet,ItemId,LocationId,RawLocationId,PlayerName,PartialAmount,TotalAmount,Received,TaxesPercent,TotalSilver,TotalTaxes,Type,UnitSilver FROM mails WHERE Deleted=0 ORDER BY Received DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Mail
	for rows.Next() {
		m, err := scanMail(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

// UnsyncedMails returns mails not yet uploaded to the user account.
func (s *Store) UnsyncedMails(limit int) ([]Mail, error) {
	if limit <= 0 {
		limit = 500
	}
	rows, err := s.db.Query(`SELECT Id,AlbionServerId,AuctionType,Deleted,IsSet,ItemId,LocationId,RawLocationId,PlayerName,PartialAmount,TotalAmount,Received,TaxesPercent,TotalSilver,TotalTaxes,Type,UnitSilver FROM mails WHERE Synced=0 AND IsSet=1 ORDER BY Received LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Mail
	for rows.Next() {
		m, err := scanMail(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

// MarkMailsSynced marks the given mail ids as uploaded.
func (s *Store) MarkMailsSynced(ids []int64) error { return s.markSynced("mails", int64sToAny(ids)) }

// UnsyncedTrades returns trades not yet uploaded to the user account.
func (s *Store) UnsyncedTrades(limit int) ([]Trade, error) {
	if limit <= 0 {
		limit = 500
	}
	rows, err := s.db.Query(`SELECT Id,AlbionServerId,Amount,DateTime,Deleted,ItemId,LocationId,Operation,PlayerName,QualityLevel,RawLocationId,SalesTaxesPercent,Type,UnitSilver FROM trades WHERE Synced=0 ORDER BY DateTime LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Trade
	for rows.Next() {
		var t Trade
		var deleted int
		var dt string
		if err := rows.Scan(&t.ID, &t.AlbionServerID, &t.Amount, &dt, &deleted, &t.ItemID, &t.LocationID, &t.Operation, &t.PlayerName, &t.QualityLevel, &t.RawLocationID, &t.SalesTaxesPercent, &t.Type, &t.UnitSilver); err != nil {
			return nil, err
		}
		t.Deleted = deleted != 0
		if parsed, e := time.Parse(time.RFC3339Nano, dt); e == nil {
			t.DateTime = parsed
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

// MarkTradesSynced marks the given trade ids as uploaded.
func (s *Store) MarkTradesSynced(ids []string) error {
	a := make([]any, len(ids))
	for i, id := range ids {
		a[i] = id
	}
	return s.markSynced("trades", a)
}

func (s *Store) markSynced(table string, ids []any) error {
	if len(ids) == 0 {
		return nil
	}
	ph := make([]string, len(ids))
	for i := range ph {
		ph[i] = "?"
	}
	_, err := s.db.Exec(`UPDATE `+table+` SET Synced=1 WHERE Id IN (`+strings.Join(ph, ",")+`)`, ids...)
	return err
}

func int64sToAny(ids []int64) []any {
	a := make([]any, len(ids))
	for i, id := range ids {
		a[i] = id
	}
	return a
}

type scanner interface {
	Scan(dest ...any) error
}

func scanMail(sc scanner) (Mail, error) {
	var m Mail
	var deleted, isSet int
	var received string
	var itemID, rawLoc, playerName sql.NullString
	// Body fields are NULL until the mail is read (InsertMailInfo writes the header only).
	var partial, total, totalSilver, totalTaxes sql.NullInt64
	var taxesPct, unitSilver sql.NullFloat64
	err := sc.Scan(&m.ID, &m.AlbionServerID, &m.AuctionType, &deleted, &isSet, &itemID, &m.LocationID, &rawLoc, &playerName, &partial, &total, &received, &taxesPct, &totalSilver, &totalTaxes, &m.Type, &unitSilver)
	if err != nil {
		return m, err
	}
	m.ItemID = itemID.String
	m.RawLocationID = rawLoc.String
	m.PlayerName = playerName.String
	m.PartialAmount = int(partial.Int64)
	m.TotalAmount = int(total.Int64)
	m.TotalSilver = totalSilver.Int64
	m.TotalTaxes = totalTaxes.Int64
	m.TaxesPercent = taxesPct.Float64
	m.UnitSilver = unitSilver.Float64
	m.Deleted = deleted != 0
	m.IsSet = isSet != 0
	if t, e := time.Parse(time.RFC3339Nano, received); e == nil {
		m.Received = t
	}
	return m, nil
}

// ---- Trades ----

// Trade is a persisted trade (instant action or filled order from mail).
type Trade struct {
	ID                string    `json:"id"`
	AlbionServerID    int       `json:"albionServerId"`
	Amount            int       `json:"amount"`
	DateTime          time.Time `json:"dateTime"`
	Deleted           bool      `json:"deleted"`
	ItemID            string    `json:"itemId"`
	LocationID        int       `json:"locationId"`
	Operation         int       `json:"operation"`
	PlayerName        string    `json:"playerName"`
	QualityLevel      int       `json:"qualityLevel"`
	RawLocationID     string    `json:"rawLocationId"`
	SalesTaxesPercent float64   `json:"salesTaxesPercent"`
	Type              int       `json:"type"`
	UnitSilver        float64   `json:"unitSilver"`
}

// InsertTrade stores a trade.
func (s *Store) InsertTrade(t Trade) error {
	_, err := s.db.Exec(`INSERT OR REPLACE INTO trades
		(Id,AlbionServerId,Amount,DateTime,Deleted,ItemId,LocationId,Operation,PlayerName,QualityLevel,RawLocationId,SalesTaxesPercent,Type,UnitSilver)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		t.ID, t.AlbionServerID, t.Amount, t.DateTime.UTC().Format(time.RFC3339Nano), b2i(t.Deleted), t.ItemID, t.LocationID, t.Operation, t.PlayerName, t.QualityLevel, t.RawLocationID, t.SalesTaxesPercent, t.Type, t.UnitSilver)
	return err
}

// ListTrades returns up to limit most-recent non-deleted trades.
func (s *Store) ListTrades(limit int) ([]Trade, error) {
	if limit <= 0 {
		limit = 1000
	}
	rows, err := s.db.Query(`SELECT Id,AlbionServerId,Amount,DateTime,Deleted,ItemId,LocationId,Operation,PlayerName,QualityLevel,RawLocationId,SalesTaxesPercent,Type,UnitSilver FROM trades WHERE Deleted=0 ORDER BY DateTime DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Trade
	for rows.Next() {
		var t Trade
		var deleted int
		var dt string
		if err := rows.Scan(&t.ID, &t.AlbionServerID, &t.Amount, &dt, &deleted, &t.ItemID, &t.LocationID, &t.Operation, &t.PlayerName, &t.QualityLevel, &t.RawLocationID, &t.SalesTaxesPercent, &t.Type, &t.UnitSilver); err != nil {
			return nil, err
		}
		t.Deleted = deleted != 0
		if parsed, e := time.Parse(time.RFC3339Nano, dt); e == nil {
			t.DateTime = parsed
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}
