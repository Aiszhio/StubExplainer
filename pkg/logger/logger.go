package logger

import (
	"log"
	"os"
)

// New returns a standard library logger used by the demo service.
func New() *log.Logger {
	return log.New(os.Stdout, "explainer: ", log.LstdFlags|log.Lmsgprefix)
}
