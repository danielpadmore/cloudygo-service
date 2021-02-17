package handlers

import (
	"net/http"

	"github.com/danielpadmore/cloudygo-service/data"
	"github.com/danielpadmore/cloudygo-service/logs"
)

// Resource contains handler data for a single resource
type Resource struct {
	logger     logs.Logger
	connection data.Connection
}

// NewResource creates a new Resource
func NewResource(logger logs.Logger, connection data.Connection) *Resource {
	return &Resource{logger, connection}
}

// ServeHTTP handles fetching all resources available
func (resource *Resource) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	resource.logger.Info(newLog("Request received at url %s", r.URL.String()))

	resources, err := resource.connection.GetResources()
	if err != nil {
		resource.logger.Error(newLog("Resources not found"))
		http.Error(rw, "Unable to find resources", http.StatusInternalServerError)
		return
	}

	data, err := resources.ToJSON()
	if err != nil {
		resource.logger.Error(newLog("Failed to parse found resources: %s", err.Error()))
		http.Error(rw, "Failed to correctly parse resources", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(data)
}
