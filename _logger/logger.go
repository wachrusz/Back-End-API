//go:build !exclude_swagger
// +build !exclude_swagger

// Package logger provides logging functionality.
package logger

import (
	"log"
	"os"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

// @Summary Initialize loggers
// @Description Initialize the info and error loggers.
// @Tags Logger
// @Router /logger [post]
func init() {
	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
