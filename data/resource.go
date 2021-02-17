package data

import "github.com/danielpadmore/cloudygo-service/model"

// GetResources returns all resources available
func (c *PostgresSQL) GetResources() (model.Resources, error) {
	resources := model.Resources{}

	err := c.db.Select(&resources, "SELECT * FROM resources WHERE available = TRUE")
	if err != nil {
		return nil, err
	}

	return resources, nil
}
