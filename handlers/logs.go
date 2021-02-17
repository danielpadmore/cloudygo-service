package handlers

import (
	"fmt"

	"github.com/danielpadmore/cloudygo-service/logs"
)

func newLog(message string, a ...interface{}) logs.LogStruct {
	return logs.NewLog("ROUTER", fmt.Sprintf(message, a...))
}
