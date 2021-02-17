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

// Lambda contains handler data for a single Lambda
type Lambda struct {
	logger     logs.Logger
	connection data.Connection
}

// NewLambda creates a new Lambda
func NewLambda(logger logs.Logger, connection data.Connection) *Lambda {
	return &Lambda{logger, connection}
}

// ServeHTTP handles fetching all resources available
func (l *Lambda) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	http.NotFound(rw, r)
}

// GetLambdas handles fetching all lambdas
func (l *Lambda) GetLambdas(userID string, rw http.ResponseWriter, r *http.Request) {
	l.logger.Info(newLog("Get lambdas request made at %s", r.URL.String()))

	lambdas, err := l.connection.GetLambdas(userID, nil)
	if err != nil {
		l.logger.Warning(newLog("Unable to find lambdas: %s", err.Error()))
		http.Error(rw, "Unable to find lambdas", http.StatusInternalServerError)
		return
	}

	data, err := lambdas.ToJSON()
	if err != nil {
		l.logger.Error(newLog("Failed to parse lambdas to JSON: %s", err.Error()))
		http.Error(rw, "Failed to correctly parse lambdas to JSON", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(data)
}

// GetLambda handles fetching a single lambda
func (l *Lambda) GetLambda(userID string, rw http.ResponseWriter, r *http.Request) {
	l.logger.Info(newLog("Get lambdas request made at %s", r.URL.String()))

	vars := mux.Vars(r)
	ID := vars["id"]

	lambdas, err := l.connection.GetLambdas(userID, &ID)
	if err != nil {
		l.logger.Warning(newLog("Error finding lambda user: %s ID: %s error: %s", userID, ID, err.Error()))
		http.Error(rw, fmt.Sprintf("Unable to find lambda %s", ID), http.StatusInternalServerError)
		return
	}

	lambda := model.Lambda{}

	if len(lambdas) > 0 {
		lambda = lambdas[0]
	} else {
		l.logger.Info(newLog("Unable to find lambda %s", ID))
		http.Error(rw, "Failed to find lambda", http.StatusNotFound)
		return
	}

	data, err := lambda.ToJSON()
	if err != nil {
		l.logger.Error(newLog("Failed to parse lambdas to JSON: %s", err.Error()))
		http.Error(rw, "Failed to correctly parse lambda to JSON", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(data)
}

// CreateLambda handles creating a new lambda
func (l *Lambda) CreateLambda(userID string, rw http.ResponseWriter, r *http.Request) {
	l.logger.Info(newLog("Create lambdas request made at %s", r.URL.String()))

	body := model.Lambda{}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		l.logger.Info(newLog("Unable to parse request body: %s", err.Error()))
		http.Error(rw, "Unable to parse request body", http.StatusBadRequest)
		return
	}

	created, err := l.connection.CreateLambda(userID, body)

	if err != nil {
		l.logger.Warning(newLog("Unable to create lambda: %s", err.Error()))
		http.Error(rw, "Unable to create lambda", http.StatusInternalServerError)
		return
	}

	data, err := created.ToJSON()
	if err != nil {
		l.logger.Error(newLog("Failed to parse lambdas to JSON: %s", err.Error()))
		http.Error(rw, "Failed to correctly parse lambda to JSON", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(data)
}

// UpdateLambda handles updating an existing lambda
func (l *Lambda) UpdateLambda(userID string, rw http.ResponseWriter, r *http.Request) {
	l.logger.Info(newLog("Update lambdas request made at %s", r.URL.String()))

	vars := mux.Vars(r)
	ID := vars["id"]

	body := model.Lambda{}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		l.logger.Info(newLog("Unable to parse request body: %s", err.Error()))
		http.Error(rw, "Unable to parse request body", http.StatusBadRequest)
		return
	}

	created, err := l.connection.UpdateLambda(userID, ID, body)

	if err != nil {
		l.logger.Warning(newLog("Unable to update lambda: %s", err.Error()))
		http.Error(rw, "Unable to update lambda", http.StatusInternalServerError)
		return
	}

	data, err := created.ToJSON()
	if err != nil {
		l.logger.Error(newLog("Failed to parse lambdas to JSON: %s", err.Error()))
		http.Error(rw, "Failed to correctly parse lambda to JSON", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(data)
}

// DeleteLambda handles deleting an existing lambda
func (l *Lambda) DeleteLambda(userID string, rw http.ResponseWriter, r *http.Request) {
	l.logger.Info(newLog("Delete lambdas request made at %s", r.URL.String()))

	vars := mux.Vars(r)
	ID := vars["id"]

	err := l.connection.DeleteLambda(userID, ID)

	if err != nil {
		l.logger.Warning(newLog("Unable to delete lambda: %s", err.Error()))
		http.Error(rw, "Unable to delete lambda", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(rw, "%s", "Lambda deleted")

}
