package model

import (
	"encoding/json"
	"io"
)

// Resource describes a type of infrastructure that cloudygo can provide
type Resource struct {
	ID        string `db:"id" json:"id,omitempty"`
	Name      string `db:"name" json:"name"`
	Type      string `db:"type" json:"type"`
	Available bool   `db:"available" json:"-"`
}

// FromJSON converts data from JSON
func (r *Resource) FromJSON(data io.Reader) error {
	de := json.NewDecoder(data)
	return de.Decode(r)
}

// ToJSON converts data to JSON
func (r *Resource) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

// Resources is a list of Resource
type Resources []Resource

// FromJSON converts data from JSON
func (r *Resources) FromJSON(data io.Reader) error {
	de := json.NewDecoder(data)
	return de.Decode(r)
}

// ToJSON converts data to JSON
func (r *Resources) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}
