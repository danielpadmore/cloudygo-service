package model

import (
	"database/sql"
	"encoding/json"
	"io"
)

// VirtualMachine is like running on a real computer... except its not!
type VirtualMachine struct {
	ID        string         `db:"id" json:"id,omitempty"`
	UserID    string         `db:"user_id" json:"-"`
	Name      string         `db:"name" json:"name"`
	Cpus      uint           `db:"cpus" json:"cpus"`
	Quantity  int            `db:"quantity" json:"quantity,omitempty"`
	CreatedAt string         `db:"created_at" json:"-"`
	UpdatedAt string         `db:"updated_at" json:"-"`
	DeletedAt sql.NullString `db:"deleted_at" json:"-"`
}

// FromJSON converts data from JSON
func (vm *VirtualMachine) FromJSON(data io.Reader) error {
	de := json.NewDecoder(data)
	return de.Decode(vm)
}

// ToJSON converts data to JSON
func (vm *VirtualMachine) ToJSON() ([]byte, error) {
	return json.Marshal(vm)
}

// VirtualMachines is a list of VirtualMachine
type VirtualMachines []VirtualMachine

// FromJSON converts data from JSON
func (r *VirtualMachines) FromJSON(data io.Reader) error {
	de := json.NewDecoder(data)
	return de.Decode(r)
}

// ToJSON converts data to JSON
func (r *VirtualMachines) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}
