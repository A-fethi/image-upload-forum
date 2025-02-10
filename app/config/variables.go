package config

import (
	"log"
	"os"
)

var (
	Logger  *log.Logger
	logFile *os.File
)

type TemplateData struct {
    IsAuthenticated bool
    Username        string
    Is404           bool
}
