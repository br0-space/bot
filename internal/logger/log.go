package logger

import (
	"os"

	"github.com/op/go-logging"
)

var Log = logging.MustGetLogger("br0fessionalbot")

var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfile} %{shortfunc} %{level}:%{color:reset} %{message}`,
)

func init() {
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)
}
