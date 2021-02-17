package data

import (
	"github.com/danielpadmore/cloudygo-service/model"
	"github.com/google/uuid"
)

// CreateVirtualMachine creates a new VirtualMachine
func (c *PostgresSQL) CreateVirtualMachine(userID string, VirtualMachine model.VirtualMachine) (model.VirtualMachine, error) {
	id := uuid.New().String()

	_, err := c.db.NamedExec(
		`INSERT INTO virtual_machines (id, user_id, name, cpus, quantity, created_at, updated_at)
		VALUES (:id, :user_id, :name, :cpus, :quantity, now(), now())`, map[string]interface{}{
			"id":       id,
			"user_id":  userID,
			"name":     VirtualMachine.Name,
			"cpus":     VirtualMachine.Cpus,
			"quantity": VirtualMachine.Quantity,
		})
	if err != nil {
		return VirtualMachine, err
	}

	VirtualMachine.ID = id

	return VirtualMachine, nil
}

// GetVirtualMachines fetches all VirtualMachines for a user with an optional filter of VirtualMachine id
func (c *PostgresSQL) GetVirtualMachines(userID string, VirtualMachineID *string) (model.VirtualMachines, error) {
	VirtualMachines := model.VirtualMachines{}

	if VirtualMachineID != nil {
		err := c.db.Select(&VirtualMachines,
			`SELECT * FROM virtual_machines WHERE user_id = $1 AND id = $2 AND deleted_at IS NULL`,
			userID, VirtualMachineID)
		if err != nil {
			return nil, err
		}
	} else {
		err := c.db.Select(&VirtualMachines,
			`SELECT * FROM virtual_machines WHERE user_id = $1 AND deleted_at IS NULL`,
			userID)
		if err != nil {
			return nil, err
		}
	}
	return VirtualMachines, nil
}

// UpdateVirtualMachine updates an existing VirtualMachine
func (c *PostgresSQL) UpdateVirtualMachine(userID string, ID string, VirtualMachine model.VirtualMachine) (model.VirtualMachine, error) {

	_, err := c.db.NamedExec(
		`UPDATE virtual_machines SET (name, cpus, quantity, updated_at) = (:name, :cpus, :quantity, now())
		WHERE id = :id AND user_id = :user_id`, map[string]interface{}{
			"id":       ID,
			"user_id":  userID,
			"name":     VirtualMachine.Name,
			"cpus":     VirtualMachine.Cpus,
			"quantity": VirtualMachine.Quantity,
		})
	if err != nil {
		return VirtualMachine, err
	}

	return VirtualMachine, nil

}

// DeleteVirtualMachine destroys an existing VirtualMachine
func (c *PostgresSQL) DeleteVirtualMachine(userID string, VirtualMachineID string) error {
	_, err := c.db.NamedExec(
		`UPDATE virtual_machines SET deleted_at = now()
		WHERE id = :id AND user_id = :user_id`, map[string]interface{}{
			"id":      VirtualMachineID,
			"user_id": userID,
		})
	if err != nil {
		return err
	}

	return nil
}
