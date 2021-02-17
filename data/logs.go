package data

import (
	"fmt"

	"github.com/danielpadmore/cloudygo-service/logs"
)

func newLog(message string, a ...interface{}) logs.LogStruct {
	return logs.NewLog("DATABASE", fmt.Sprintf(message, a...))
}
