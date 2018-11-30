package logger

import (
	"log"
	"os"
)

var (
	Info = log.New(os.Stdout, "[INFO] ", log.LstdFlags)
)
