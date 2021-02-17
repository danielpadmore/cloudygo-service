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

// NoSQLDatabase contains handler data for a single NoSQLDatabase
type NoSQLDatabase struct {
	logger     logs.Logger
	connection data.Connection
}

// NewNoSQLDatabase creates a new NoSQLDatabase
func NewNoSQLDatabase(logger logs.Logger, connection data.Connection) *NoSQLDatabase {
	return &NoSQLDatabase{logger, connection}
}

// ServeHTTP handles fetching all resources available
func (l *NoSQLDatabase) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	http.NotFound(rw, r)
}

// GetNoSQLDatabases handles fetching all NoSQLDatabases
func (l *NoSQLDatabase) GetNoSQLDatabases(userID string, rw http.ResponseWriter, r *http.Request) {
	l.logger.Info(newLog("Get NoSQL databases request made at %s", r.URL.String()))

	res, err := l.connection.GetNoSQLDatabases(userID, nil)
	if err != nil {
		l.logger.Warning(newLog("Unable to find NoSQL databases: %s", err.Error()))
		http.Error(rw, "Unable to find NoSQL databases", http.StatusInternalServerError)
		return
	}

	data, err := res.ToJSON()
	if err != nil {
		l.logger.Error(newLog("Failed to parse NoSQL databases to JSON: %s", err.Error()))
		http.Error(rw, "Failed to correctly parse NoSQL databases to JSON", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(data)
}

// GetNoSQLDatabase handles fetching a single NoSQLDatabase
func (l *NoSQLDatabase) GetNoSQLDatabase(userID string, rw http.ResponseWriter, r *http.Request) {
	l.logger.Info(newLog("Get NoSQL databases request made at %s", r.URL.String()))

	vars := mux.Vars(r)
	ID := vars["id"]

	res, err := l.connection.GetNoSQLDatabases(userID, &ID)
	if err != nil {
		l.logger.Warning(newLog("Error finding NoSQL database user: %s ID: %s error: %s", userID, ID, err.Error()))
		http.Error(rw, fmt.Sprintf("Unable to find NoSQL database %s", ID), http.StatusInternalServerError)
		return
	}

	db := model.NoSQLDatabase{}

	if len(res) > 0 {
		db = res[0]
	} else {
		l.logger.Info(newLog("Unable to find NoSQL database %s", ID))
		http.Error(rw, "Failed to find NoSQL database", http.StatusNotFound)
		return
	}

	data, err := db.ToJSON()
	if err != nil {
		l.logger.Error(newLog("Failed to parse NoSQL database to JSON: %s", err.Error()))
		http.Error(rw, "Failed to correctly parse NoSQL database to JSON", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(data)
}

// CreateNoSQLDatabase handles creating a new NoSQLDatabase
func (l *NoSQLDatabase) CreateNoSQLDatabase(userID string, rw http.ResponseWriter, r *http.Request) {
	l.logger.Info(newLog("Create NoSQL database request made at %s", r.URL.String()))

	body := model.NoSQLDatabase{}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		l.logger.Info(newLog("Unable to parse request body: %s", err.Error()))
		http.Error(rw, "Unable to parse request body", http.StatusBadRequest)
		return
	}

	created, err := l.connection.CreateNoSQLDatabase(userID, body)

	if err != nil {
		l.logger.Warning(newLog("Unable to create NoSQL database: %s", err.Error()))
		http.Error(rw, "Unable to create NoSQL database", http.StatusInternalServerError)
		return
	}

	data, err := created.ToJSON()
	if err != nil {
		l.logger.Error(newLog("Failed to parse NoSQL database to JSON: %s", err.Error()))
		http.Error(rw, "Failed to correctly parse NoSQL database to JSON", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(data)
}

// UpdateNoSQLDatabase handles updating an existing NoSQLDatabase
func (l *NoSQLDatabase) UpdateNoSQLDatabase(userID string, rw http.ResponseWriter, r *http.Request) {
	l.logger.Info(newLog("Update NoSQL database request made at %s", r.URL.String()))

	vars := mux.Vars(r)
	ID := vars["id"]

	body := model.NoSQLDatabase{}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		l.logger.Info(newLog("Unable to parse request body: %s", err.Error()))
		http.Error(rw, "Unable to parse request body", http.StatusBadRequest)
		return
	}

	created, err := l.connection.UpdateNoSQLDatabase(userID, ID, body)

	if err != nil {
		l.logger.Warning(newLog("Unable to update NoSQL database: %s", err.Error()))
		http.Error(rw, "Unable to update NoSQL database", http.StatusInternalServerError)
		return
	}

	data, err := created.ToJSON()
	if err != nil {
		l.logger.Error(newLog("Failed to parse NoSQL database to JSON: %s", err.Error()))
		http.Error(rw, "Failed to correctly parse NoSQL database to JSON", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(data)
}

// DeleteNoSQLDatabase handles deleting an existing NoSQLDatabase
func (l *NoSQLDatabase) DeleteNoSQLDatabase(userID string, rw http.ResponseWriter, r *http.Request) {
	l.logger.Info(newLog("Delete NoSQL database request made at %s", r.URL.String()))

	vars := mux.Vars(r)
	ID := vars["id"]

	err := l.connection.DeleteNoSQLDatabase(userID, ID)

	if err != nil {
		l.logger.Warning(newLog("Unable to delete NoSQL database: %s", err.Error()))
		http.Error(rw, "Unable to delete NoSQL database", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(rw, "%s", "NoSQL database deleted")

}
