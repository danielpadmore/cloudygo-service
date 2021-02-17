package data

import (
	"github.com/danielpadmore/cloudygo-service/model"
	"github.com/google/uuid"
)

// CreateSQLDatabase creates a new SQLDatabase
func (c *PostgresSQL) CreateSQLDatabase(userID string, SQLDatabase model.SQLDatabase) (model.SQLDatabase, error) {
	id := uuid.New().String()

	_, err := c.db.NamedExec(
		`INSERT INTO sql_databases (id, user_id, name, username, password, quantity, created_at, updated_at)
		VALUES (:id, :user_id, :name, :username, :password, :quantity, now(), now())`, map[string]interface{}{
			"id":       id,
			"user_id":  userID,
			"name":     SQLDatabase.Name,
			"username": SQLDatabase.Username,
			"password": SQLDatabase.Password,
			"quantity": SQLDatabase.Quantity,
		})
	if err != nil {
		return SQLDatabase, err
	}

	SQLDatabase.ID = id

	return SQLDatabase, nil
}

// GetSQLDatabases fetches all SQLDatabases for a user with an optional filter of SQLDatabase id
func (c *PostgresSQL) GetSQLDatabases(userID string, SQLDatabaseID *string) (model.SQLDatabases, error) {
	SQLDatabases := model.SQLDatabases{}

	if SQLDatabaseID != nil {
		err := c.db.Select(&SQLDatabases,
			`SELECT * FROM sql_databases WHERE user_id = $1 AND id = $2 AND deleted_at IS NULL`,
			userID, SQLDatabaseID)
		if err != nil {
			return nil, err
		}
	} else {
		err := c.db.Select(&SQLDatabases,
			`SELECT * FROM sql_databases WHERE user_id = $1 AND deleted_at IS NULL`,
			userID)
		if err != nil {
			return nil, err
		}
	}
	return SQLDatabases, nil
}

// UpdateSQLDatabase updates an existing SQLDatabase
func (c *PostgresSQL) UpdateSQLDatabase(userID string, ID string, SQLDatabase model.SQLDatabase) (model.SQLDatabase, error) {

	_, err := c.db.NamedExec(
		`UPDATE sql_databases SET (name, username, password, quantity, updated_at) = (:name, :username, :password, :quantity, now())
		WHERE id = :id AND user_id = :user_id`, map[string]interface{}{
			"id":       ID,
			"user_id":  userID,
			"name":     SQLDatabase.Name,
			"username": SQLDatabase.Username,
			"password": SQLDatabase.Password,
			"quantity": SQLDatabase.Quantity,
		})
	if err != nil {
		return SQLDatabase, err
	}

	return SQLDatabase, nil

}

// DeleteSQLDatabase destroys an existing SQLDatabase
func (c *PostgresSQL) DeleteSQLDatabase(userID string, SQLDatabaseID string) error {
	_, err := c.db.NamedExec(
		`UPDATE sql_databases SET deleted_at = now()
		WHERE id = :id AND user_id = :user_id`, map[string]interface{}{
			"id":      SQLDatabaseID,
			"user_id": userID,
		})
	if err != nil {
		return err
	}

	return nil
}
