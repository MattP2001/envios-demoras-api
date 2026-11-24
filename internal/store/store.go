package store

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

type Event struct {
	ID          int64     `json:"id"`
	ShipmentID  string    `json:"shipment_id"`
	Status      string    `json:"status"`
	EventDate   time.Time `json:"event_date"`
	ReceivedAt  time.Time `json:"received_at"`
}

type Shipment struct {
	ShipmentID    string     `json:"shipment_id"`
	Status        string     `json:"status"`
	PromisedDate  *time.Time `json:"promised_date"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type Store struct {
	db *sql.DB
}

func New(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS shipments (
  shipment_id TEXT PRIMARY KEY,
  status TEXT NOT NULL,
  promised_date TEXT,
  updated_at TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS shipment_events (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  shipment_id TEXT NOT NULL,
  status TEXT NOT NULL,
  event_date TEXT NOT NULL,
  received_at TEXT NOT NULL
);`)
	if err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

// SaveEvent guarda el evento y actualiza el estado del envio con el ultimo que llego.
func (s *Store) SaveEvent(shipmentID, status string, eventDate time.Time, promised *time.Time) error {
	now := time.Now().UTC()
	_, err := s.db.Exec(`INSERT INTO shipment_events (shipment_id, status, event_date, received_at) VALUES (?, ?, ?, ?)`,
		shipmentID, status, eventDate.Format(time.RFC3339), now.Format(time.RFC3339))
	if err != nil {
		return err
	}
	var p any
	if promised != nil {
		p = promised.Format(time.RFC3339)
	}
	_, err = s.db.Exec(`INSERT INTO shipments (shipment_id, status, promised_date, updated_at) VALUES (?, ?, ?, ?)
ON CONFLICT(shipment_id) DO UPDATE SET status = excluded.status,
  promised_date = COALESCE(excluded.promised_date, shipments.promised_date),
  updated_at = excluded.updated_at`,
		shipmentID, status, p, now.Format(time.RFC3339))
	return err
}

func (s *Store) GetShipment(id string) (*Shipment, error) {
	var sh Shipment
	var promised sql.NullString
	var updated string
	err := s.db.QueryRow(`SELECT shipment_id, status, promised_date, updated_at FROM shipments WHERE shipment_id = ?`, id).
		Scan(&sh.ShipmentID, &sh.Status, &promised, &updated)
	if err != nil {
		return nil, err
	}
	if promised.Valid {
		t, _ := time.Parse(time.RFC3339, promised.String)
		sh.PromisedDate = &t
	}
	sh.UpdatedAt, _ = time.Parse(time.RFC3339, updated)
	return &sh, nil
}

func (s *Store) GetEvents(id string) ([]Event, error) {
	rows, err := s.db.Query(`SELECT id, shipment_id, status, event_date, received_at FROM shipment_events WHERE shipment_id = ? ORDER BY id`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Event
	for rows.Next() {
		var e Event
		var ed, ra string
		if err := rows.Scan(&e.ID, &e.ShipmentID, &e.Status, &ed, &ra); err != nil {
			return nil, err
		}
		e.EventDate, _ = time.Parse(time.RFC3339, ed)
		e.ReceivedAt, _ = time.Parse(time.RFC3339, ra)
		out = append(out, e)
	}
	return out, nil
}
