package data

import (
	"errors"

	"github.com/danielpadmore/cloudygo-service/model"
	"github.com/google/uuid"
)

// CreateUser creates a new user in the users table
func (c *PostgresSQL) CreateUser(username string, password string) (model.User, error) {
	c.logger.Verbose(newLog("Create user called"))
	user := model.User{}

	id := uuid.New()

	rows, err := c.db.NamedQuery(
		`INSERT INTO users (id, username, password, created_at, updated_at) 
		VALUES(:id, :username, crypt(:password, gen_salt('bf')), now(), now()) 
		RETURNING id, username;`, map[string]interface{}{
			"id":       id.String(),
			"username": username,
			"password": password,
		})
	if err != nil {
		c.logger.Info(newLog("Error creating user: %s", err.Error()))
		return user, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.StructScan(&user)
		if err != nil {
			c.logger.Error(newLog("Error parsing user record: %s", err.Error()))
			return user, err
		}
	}

	c.logger.Info(newLog("Created user %s", user.ID))
	return user, err
}

// AuthenticateUser ensures username and password match and returns result
func (c *PostgresSQL) AuthenticateUser(username string, password string) (model.User, error) {
	c.logger.Verbose(newLog("Authenticate user called"))
	users := []model.User{}

	err := c.db.Select(&users,
		`SELECT id, username FROM users WHERE username = $1 AND password = crypt($2, password);`,
		username,
		password,
	)

	if err != nil {
		c.logger.Info(newLog("Error authenticating user: %s", err.Error()))
		return model.User{}, err
	}

	if len(users) <= 0 {
		c.logger.Info(newLog("User not found"))
		return model.User{}, errors.New("User not found")
	}

	c.logger.Info(newLog("User %s found", users[0].Username))
	return users[0], nil
}
