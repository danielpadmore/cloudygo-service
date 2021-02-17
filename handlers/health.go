package handlers

import (
	"net/http"

	"github.com/danielpadmore/cloudygo-service/data"
	"github.com/danielpadmore/cloudygo-service/logs"
)

// Health contains database connection data
type Health struct {
	logger logs.Logger
	db     data.Connection
}

// NewHealth creates a new Health instance
func NewHealth(logger logs.Logger, db data.Connection) *Health {
	return &Health{logger, db}
}

func (h *Health) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	_, err := h.db.IsConnected()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		h.logger.Error(newLog("Error with database connection: %s", err.Error()))
	}

	h.logger.Error(newLog("Health check ok"))
}
