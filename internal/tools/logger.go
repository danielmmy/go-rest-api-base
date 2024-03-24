package tools

import (
	"log"
	"os"
)

var InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
var WarningLogger = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
var ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
