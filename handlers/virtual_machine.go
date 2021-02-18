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

// VirtualMachine contains handler data for a single VirtualMachine
type VirtualMachine struct {
	logger     logs.Logger
	connection data.Connection
}

type createVirtualMachineRequestBody struct {
	Name     string `json:"name" validate:"required,min=5,max=200"`
	Cpus     uint   `json:"cpus" validate:"required,gte=1,lte=64"`
	Quantity int    `json:"quantity" validate:"required,gte=1,lte=500"`
}

// NewVirtualMachine creates a new VirtualMachine
func NewVirtualMachine(logger logs.Logger, connection data.Connection) *VirtualMachine {
	return &VirtualMachine{logger, connection}
}

// ServeHTTP handles fetching all resources available
func (l *VirtualMachine) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	http.NotFound(rw, r)
}

// GetVirtualMachines handles fetching all VirtualMachines
func (l *VirtualMachine) GetVirtualMachines(userID string, rw http.ResponseWriter, r *http.Request) {
	l.logger.Info(newLog("Get virtual machines request made at %s", r.URL.String()))

	res, err := l.connection.GetVirtualMachines(userID, nil)
	if err != nil {
		l.logger.Warning(newLog("Unable to find virtual machines: %s", err.Error()))
		http.Error(rw, "Unable to find virtual machines", http.StatusInternalServerError)
		return
	}

	data, err := res.ToJSON()
	if err != nil {
		l.logger.Error(newLog("Failed to parse virtual machines to JSON: %s", err.Error()))
		http.Error(rw, "Failed to correctly parse virtual machines to JSON", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(data)
}

// GetVirtualMachine handles fetching a single VirtualMachine
func (l *VirtualMachine) GetVirtualMachine(userID string, rw http.ResponseWriter, r *http.Request) {
	l.logger.Info(newLog("Get virtual machines request made at %s", r.URL.String()))

	vars := mux.Vars(r)
	ID := vars["id"]

	res, err := l.connection.GetVirtualMachines(userID, &ID)
	if err != nil {
		l.logger.Warning(newLog("Error finding virtual machine user: %s ID: %s error: %s", userID, ID, err.Error()))
		http.Error(rw, fmt.Sprintf("Unable to find virtual machine %s", ID), http.StatusInternalServerError)
		return
	}

	vm := model.VirtualMachine{}

	if len(res) > 0 {
		vm = res[0]
	} else {
		l.logger.Info(newLog("Unable to find virtual machine %s", ID))
		http.Error(rw, "Failed to find virtual machine", http.StatusNotFound)
		return
	}

	data, err := vm.ToJSON()
	if err != nil {
		l.logger.Error(newLog("Failed to parse virtual machine to JSON: %s", err.Error()))
		http.Error(rw, "Failed to correctly parse virtual machine to JSON", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(data)
}

// CreateVirtualMachine handles creating a new VirtualMachine
func (l *VirtualMachine) CreateVirtualMachine(userID string, rw http.ResponseWriter, r *http.Request) {
	l.logger.Info(newLog("Create virtual machine request made at %s", r.URL.String()))

	body := model.VirtualMachine{}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		l.logger.Info(newLog("Unable to parse request body: %s", err.Error()))
		http.Error(rw, "Unable to parse request body", http.StatusBadRequest)
		return
	}

	created, err := l.connection.CreateVirtualMachine(userID, body)

	if err != nil {
		l.logger.Warning(newLog("Unable to create virtual machine: %s", err.Error()))
		http.Error(rw, "Unable to create virtual machine", http.StatusInternalServerError)
		return
	}

	data, err := created.ToJSON()
	if err != nil {
		l.logger.Error(newLog("Failed to parse virtual machine to JSON: %s", err.Error()))
		http.Error(rw, "Failed to correctly parse virtual machine to JSON", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(data)
}

// UpdateVirtualMachine handles updating an existing VirtualMachine
func (l *VirtualMachine) UpdateVirtualMachine(userID string, rw http.ResponseWriter, r *http.Request) {
	l.logger.Info(newLog("Update virtual machine request made at %s", r.URL.String()))

	vars := mux.Vars(r)
	ID := vars["id"]

	body := model.VirtualMachine{}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		l.logger.Info(newLog("Unable to parse request body: %s", err.Error()))
		http.Error(rw, "Unable to parse request body", http.StatusBadRequest)
		return
	}

	created, err := l.connection.UpdateVirtualMachine(userID, ID, body)

	if err != nil {
		l.logger.Warning(newLog("Unable to update virtual machine: %s", err.Error()))
		http.Error(rw, "Unable to update virtual machine", http.StatusInternalServerError)
		return
	}

	data, err := created.ToJSON()
	if err != nil {
		l.logger.Error(newLog("Failed to parse virtual machine to JSON: %s", err.Error()))
		http.Error(rw, "Failed to correctly parse virtual machine to JSON", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(data)
}

// DeleteVirtualMachine handles deleting an existing VirtualMachine
func (l *VirtualMachine) DeleteVirtualMachine(userID string, rw http.ResponseWriter, r *http.Request) {
	l.logger.Info(newLog("Delete virtual machine request made at %s", r.URL.String()))

	vars := mux.Vars(r)
	ID := vars["id"]

	err := l.connection.DeleteVirtualMachine(userID, ID)

	if err != nil {
		l.logger.Warning(newLog("Unable to delete virtual machine: %s", err.Error()))
		http.Error(rw, "Unable to delete virtual machine", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(rw, "%s", "virtual machine deleted")

}
