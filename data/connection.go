package data

import (
	"github.com/danielpadmore/cloudygo-service/logs"
	"github.com/danielpadmore/cloudygo-service/model"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Postgres driver
)

// Connection interface controls methods for database access patterns
type Connection interface {
	IsConnected() (bool, error)
	GetResources() (model.Resources, error)
	CreateUser(string, string) (model.User, error)
	AuthenticateUser(string, string) (model.User, error)
	CreateLambda(string, model.Lambda) (model.Lambda, error)
	GetLambdas(string, *string) (model.Lambdas, error)
	UpdateLambda(string, string, model.Lambda) (model.Lambda, error)
	DeleteLambda(string, string) error
	CreateVirtualMachine(string, model.VirtualMachine) (model.VirtualMachine, error)
	GetVirtualMachines(string, *string) (model.VirtualMachines, error)
	UpdateVirtualMachine(string, string, model.VirtualMachine) (model.VirtualMachine, error)
	DeleteVirtualMachine(string, string) error
	CreateSQLDatabase(string, model.SQLDatabase) (model.SQLDatabase, error)
	GetSQLDatabases(string, *string) (model.SQLDatabases, error)
	UpdateSQLDatabase(string, string, model.SQLDatabase) (model.SQLDatabase, error)
	DeleteSQLDatabase(string, string) error
	CreateNoSQLDatabase(string, model.NoSQLDatabase) (model.NoSQLDatabase, error)
	GetNoSQLDatabases(string, *string) (model.NoSQLDatabases, error)
	UpdateNoSQLDatabase(string, string, model.NoSQLDatabase) (model.NoSQLDatabase, error)
	DeleteNoSQLDatabase(string, string) error
}

// PostgresSQL contains database connection data
type PostgresSQL struct {
	logger logs.Logger
	db     *sqlx.DB
}

// New creates a new connection to the database
func New(logger logs.Logger, connection string) (Connection, error) {
	db, err := sqlx.Connect("postgres", connection)
	if err != nil {
		return nil, err
	}

	return &PostgresSQL{logger, db}, nil
}

// IsConnected checks the connection to the database and returns an error if not connected
func (c *PostgresSQL) IsConnected() (bool, error) {
	err := c.db.Ping()
	if err != nil {
		return false, err
	}

	return true, nil
}
