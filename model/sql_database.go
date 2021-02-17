package model

import (
	"database/sql"
	"encoding/json"
	"io"
)

// SQLDatabase is a old-school DB boi
type SQLDatabase struct {
	ID        string         `db:"id" json:"id,omitempty"`
	UserID    string         `db:"user_id" json:"-"`
	Name      string         `db:"name" json:"name"`
	Username  string         `db:"username" json:"username"`
	Password  string         `db:"password" json:"-"`
	Quantity  int            `db:"quantity" json:"quantity,omitempty"`
	CreatedAt string         `db:"created_at" json:"-"`
	UpdatedAt string         `db:"updated_at" json:"-"`
	DeletedAt sql.NullString `db:"deleted_at" json:"-"`
}

// FromJSON converts data from JSON
func (db *SQLDatabase) FromJSON(data io.Reader) error {
	de := json.NewDecoder(data)
	return de.Decode(db)
}

// ToJSON converts data to JSON
func (db *SQLDatabase) ToJSON() ([]byte, error) {
	return json.Marshal(db)
}

// SQLDatabases is a list of SQLDatabase
type SQLDatabases []SQLDatabase

// FromJSON converts data from JSON
func (r *SQLDatabases) FromJSON(data io.Reader) error {
	de := json.NewDecoder(data)
	return de.Decode(r)
}

// ToJSON converts data to JSON
func (r *SQLDatabases) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}
