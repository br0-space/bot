package logger

import (
	"os"

	"github.com/op/go-logging"
)

var liveLogger *logging.Logger

func NewLiveLogger() *logging.Logger {
	if liveLogger == nil {
		backend := logging.NewLogBackend(os.Stderr, "", 0)
		format := logging.MustStringFormatter(`%{color}%{time:15:04:05.000} %{shortfile} %{shortfunc} %{level}:%{color:reset} %{message}`)
		backendFormatter := logging.NewBackendFormatter(backend, format)
		logging.SetBackend(backendFormatter)
		liveLogger = logging.MustGetLogger("br0fessionalbot")
	}
	return liveLogger
}