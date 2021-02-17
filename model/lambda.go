package model

import (
	"database/sql"
	"encoding/json"
	"io"
)

// Lambda is a temporary server which does its business then cleans up after itself... proper neat!
type Lambda struct {
	ID              string         `db:"id" json:"id,omitempty"`
	UserID          string         `db:"user_id" json:"-"`
	Name            string         `db:"name" json:"name"`
	ConcurrentLimit uint           `db:"concurrent_limit" json:"concurrent_limit"`
	CreatedAt       string         `db:"created_at" json:"-"`
	UpdatedAt       string         `db:"updated_at" json:"-"`
	DeletedAt       sql.NullString `db:"deleted_at" json:"-"`
}

// FromJSON converts data from JSON
func (l *Lambda) FromJSON(data io.Reader) error {
	de := json.NewDecoder(data)
	return de.Decode(l)
}

// ToJSON converts data to JSON
func (l *Lambda) ToJSON() ([]byte, error) {
	return json.Marshal(l)
}

// Lambdas is a list of Lambda
type Lambdas []Lambda

// FromJSON converts data from JSON
func (l *Lambdas) FromJSON(data io.Reader) error {
	de := json.NewDecoder(data)
	return de.Decode(l)
}

// ToJSON converts data to JSON
func (l *Lambdas) ToJSON() ([]byte, error) {
	return json.Marshal(l)
}
