package data

import (
	"github.com/danielpadmore/cloudygo-service/model"
	"github.com/google/uuid"
)

// CreateNoSQLDatabase creates a new NoSQLDatabase
func (c *PostgresSQL) CreateNoSQLDatabase(userID string, NoSQLDatabase model.NoSQLDatabase) (model.NoSQLDatabase, error) {
	id := uuid.New().String()

	_, err := c.db.NamedExec(
		`INSERT INTO nosql_databases (id, user_id, name, shards, created_at, updated_at)
		VALUES (:id, :user_id, :name, :shards, now(), now())`, map[string]interface{}{
			"id":      id,
			"user_id": userID,
			"name":    NoSQLDatabase.Name,
			"shards":  NoSQLDatabase.Shards,
		})
	if err != nil {
		return NoSQLDatabase, err
	}

	NoSQLDatabase.ID = id

	return NoSQLDatabase, nil
}

// GetNoSQLDatabases fetches all NoSQLDatabases for a user with an optional filter of NoSQLDatabase id
func (c *PostgresSQL) GetNoSQLDatabases(userID string, NoSQLDatabaseID *string) (model.NoSQLDatabases, error) {
	NoSQLDatabases := model.NoSQLDatabases{}

	if NoSQLDatabaseID != nil {
		err := c.db.Select(&NoSQLDatabases,
			`SELECT * FROM nosql_databases WHERE user_id = $1 AND id = $2 AND deleted_at IS NULL`,
			userID, NoSQLDatabaseID)
		if err != nil {
			return nil, err
		}
	} else {
		err := c.db.Select(&NoSQLDatabases,
			`SELECT * FROM nosql_databases WHERE user_id = $1 AND deleted_at IS NULL`,
			userID)
		if err != nil {
			return nil, err
		}
	}
	return NoSQLDatabases, nil
}

// UpdateNoSQLDatabase updates an existing NoSQLDatabase
func (c *PostgresSQL) UpdateNoSQLDatabase(userID string, ID string, NoSQLDatabase model.NoSQLDatabase) (model.NoSQLDatabase, error) {

	_, err := c.db.NamedExec(
		`UPDATE nosql_databases SET (name, shards, updated_at) = (:name, :shards, now())
		WHERE id = :id AND user_id = :user_id`, map[string]interface{}{
			"id":      ID,
			"user_id": userID,
			"name":    NoSQLDatabase.Name,
			"shards":  NoSQLDatabase.Shards,
		})
	if err != nil {
		return NoSQLDatabase, err
	}

	return NoSQLDatabase, nil

}

// DeleteNoSQLDatabase destroys an existing NoSQLDatabase
func (c *PostgresSQL) DeleteNoSQLDatabase(userID string, NoSQLDatabaseID string) error {
	_, err := c.db.NamedExec(
		`UPDATE nosql_databases SET deleted_at = now()
		WHERE id = :id AND user_id = :user_id`, map[string]interface{}{
			"id":      NoSQLDatabaseID,
			"user_id": userID,
		})
	if err != nil {
		return err
	}

	return nil
}
