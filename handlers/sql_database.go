package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/danielpadmore/cloudygo-service/data"
	"github.com/danielpadmore/cloudygo-service/logs"
	"github.com/danielpadmore/cloudygo-service/model"
	"github.com/gorilla/mux"
)

// SQLDatabase contains handler data for a single SQLDatabase
type SQLDatabase struct {
	logger     logs.Logger
	connection data.Connection
}

// NewSQLDatabase creates a new SQLDatabase
func NewSQLDatabase(logger logs.Logger, connection data.Connection) *SQLDatabase {
	return &SQLDatabase{logger, connection}
}

// ServeHTTP handles fetching all resources available
func (l *SQLDatabase) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	http.NotFound(rw, r)
}

// GetSQLDatabases handles fetching all SQLDatabases
func (l *SQLDatabase) GetSQLDatabases(userID string, rw http.ResponseWriter, r *http.Request) {
	l.logger.Info(newLog("Get SQL databases request made at %s", r.URL.String()))

	res, err := l.connection.GetSQLDatabases(userID, nil)
	if err != nil {
		l.logger.Warning(newLog("Unable to find SQL databases: %s", err.Error()))
		http.Error(rw, "Unable to find SQL databases", http.StatusInternalServerError)
		return
	}

	data, err := res.ToJSON()
	if err != nil {
		l.logger.Error(newLog("Failed to parse SQL databases to JSON: %s", err.Error()))
		http.Error(rw, "Failed to correctly parse SQL databases to JSON", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(data)
}

// GetSQLDatabase handles fetching a single SQLDatabase
func (l *SQLDatabase) GetSQLDatabase(userID string, rw http.ResponseWriter, r *http.Request) {
	l.logger.Info(newLog("Get SQL databases request made at %s", r.URL.String()))

	vars := mux.Vars(r)
	ID := vars["id"]

	res, err := l.connection.GetSQLDatabases(userID, &ID)
	if err != nil {
		l.logger.Warning(newLog("Error finding SQL database user: %s ID: %s error: %s", userID, ID, err.Error()))
		http.Error(rw, fmt.Sprintf("Unable to find SQL database %s", ID), http.StatusInternalServerError)
		return
	}

	db := model.SQLDatabase{}

	if len(res) > 0 {
		db = res[0]
	} else {
		l.logger.Info(newLog("Unable to find SQL database %s", ID))
		http.Error(rw, "Failed to find SQL database", http.StatusNotFound)
		return
	}

	data, err := db.ToJSON()
	if err != nil {
		l.logger.Error(newLog("Failed to parse SQL database to JSON: %s", err.Error()))
		http.Error(rw, "Failed to correctly parse SQL database to JSON", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(data)
}

// CreateSQLDatabase handles creating a new SQLDatabase
func (l *SQLDatabase) CreateSQLDatabase(userID string, rw http.ResponseWriter, r *http.Request) {
	l.logger.Info(newLog("Create SQL database request made at %s", r.URL.String()))

	body := model.SQLDatabase{}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		l.logger.Info(newLog("Unable to parse request body: %s", err.Error()))
		http.Error(rw, "Unable to parse request body", http.StatusBadRequest)
		return
	}

	created, err := l.connection.CreateSQLDatabase(userID, body)

	if err != nil {
		l.logger.Warning(newLog("Unable to create SQL database: %s", err.Error()))
		http.Error(rw, "Unable to create SQL database", http.StatusInternalServerError)
		return
	}

	data, err := created.ToJSON()
	if err != nil {
		l.logger.Error(newLog("Failed to parse SQL database to JSON: %s", err.Error()))
		http.Error(rw, "Failed to correctly parse SQL database to JSON", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(data)
}

// UpdateSQLDatabase handles updating an existing SQLDatabase
func (l *SQLDatabase) UpdateSQLDatabase(userID string, rw http.ResponseWriter, r *http.Request) {
	l.logger.Info(newLog("Update SQL database request made at %s", r.URL.String()))

	vars := mux.Vars(r)
	ID := vars["id"]

	body := model.SQLDatabase{}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		l.logger.Info(newLog("Unable to parse request body: %s", err.Error()))
		http.Error(rw, "Unable to parse request body", http.StatusBadRequest)
		return
	}

	created, err := l.connection.UpdateSQLDatabase(userID, ID, body)

	if err != nil {
		l.logger.Warning(newLog("Unable to update SQL database: %s", err.Error()))
		http.Error(rw, "Unable to update SQL database", http.StatusInternalServerError)
		return
	}

	data, err := created.ToJSON()
	if err != nil {
		l.logger.Error(newLog("Failed to parse SQL database to JSON: %s", err.Error()))
		http.Error(rw, "Failed to correctly parse SQL database to JSON", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(data)
}

// DeleteSQLDatabase handles deleting an existing SQLDatabase
func (l *SQLDatabase) DeleteSQLDatabase(userID string, rw http.ResponseWriter, r *http.Request) {
	l.logger.Info(newLog("Delete SQL database request made at %s", r.URL.String()))

	vars := mux.Vars(r)
	ID := vars["id"]

	err := l.connection.DeleteSQLDatabase(userID, ID)

	if err != nil {
		l.logger.Warning(newLog("Unable to delete SQL database: %s", err.Error()))
		http.Error(rw, "Unable to delete SQL database", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(rw, "%s", "SQL database deleted")

}
