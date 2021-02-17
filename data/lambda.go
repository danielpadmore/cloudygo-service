package data

import (
	"github.com/danielpadmore/cloudygo-service/model"
	"github.com/google/uuid"
)

// CreateLambda creates a new Lambda
func (c *PostgresSQL) CreateLambda(userID string, lambda model.Lambda) (model.Lambda, error) {
	id := uuid.New().String()

	_, err := c.db.NamedExec(
		`INSERT INTO lambdas (id, user_id, name, concurrent_limit, created_at, updated_at)
		VALUES (:id, :user_id, :name, :concurrent_limit, now(), now())`, map[string]interface{}{
			"id":               id,
			"user_id":          userID,
			"name":             lambda.Name,
			"concurrent_limit": lambda.ConcurrentLimit,
		})
	if err != nil {
		return lambda, err
	}

	lambda.ID = id

	return lambda, nil
}

// GetLambdas fetches all lambdas for a user with an optional filter of lambda id
func (c *PostgresSQL) GetLambdas(userID string, lambdaID *string) (model.Lambdas, error) {
	lambdas := model.Lambdas{}

	if lambdaID != nil {
		err := c.db.Select(&lambdas,
			`SELECT * FROM lambdas WHERE user_id = $1 AND id = $2 AND deleted_at IS NULL`,
			userID, lambdaID)
		if err != nil {
			return nil, err
		}
	} else {
		err := c.db.Select(&lambdas,
			`SELECT * FROM lambdas WHERE user_id = $1 AND deleted_at IS NULL`,
			userID)
		if err != nil {
			return nil, err
		}
	}
	return lambdas, nil
}

// UpdateLambda updates an existing lambda
func (c *PostgresSQL) UpdateLambda(userID string, ID string, lambda model.Lambda) (model.Lambda, error) {

	_, err := c.db.NamedExec(
		`UPDATE lambdas SET (name, concurrent_limit, updated_at) = (:name, :concurrent_limit, now())
		WHERE id = :id AND user_id = :user_id`, map[string]interface{}{
			"id":               ID,
			"user_id":          userID,
			"name":             lambda.Name,
			"concurrent_limit": lambda.ConcurrentLimit,
		})
	if err != nil {
		return lambda, err
	}

	return lambda, nil

}

// DeleteLambda destroys an existing lambda
func (c *PostgresSQL) DeleteLambda(userID string, lambdaID string) error {
	_, err := c.db.NamedExec(
		`UPDATE lambdas SET deleted_at = now()
		WHERE id = :id AND user_id = :user_id`, map[string]interface{}{
			"id":      lambdaID,
			"user_id": userID,
		})
	if err != nil {
		return err
	}

	return nil
}
