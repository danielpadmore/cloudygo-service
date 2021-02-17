package model

import (
	"database/sql"
	"encoding/json"
	"io"
)

// NoSQLDatabase is an edgy new way to store data
type NoSQLDatabase struct {
	ID        string         `db:"id" json:"id,omitempty"`
	UserID    string         `db:"user_id" json:"-"`
	Name      string         `db:"name" json:"name"`
	Shards    uint           `db:"shards" json:"shards"`
	CreatedAt string         `db:"created_at" json:"-"`
	UpdatedAt string         `db:"updated_at" json:"-"`
	DeletedAt sql.NullString `db:"deleted_at" json:"-"`
}

// FromJSON converts data from JSON
func (db *NoSQLDatabase) FromJSON(data io.Reader) error {
	de := json.NewDecoder(data)
	return de.Decode(db)
}

// ToJSON converts data to JSON
func (db *NoSQLDatabase) ToJSON() ([]byte, error) {
	return json.Marshal(db)
}

// NoSQLDatabases is a list of NoSQLDatabase
type NoSQLDatabases []NoSQLDatabase

// FromJSON converts data from JSON
func (r *NoSQLDatabases) FromJSON(data io.Reader) error {
	de := json.NewDecoder(data)
	return de.Decode(r)
}

// ToJSON converts data to JSON
func (r *NoSQLDatabases) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}
