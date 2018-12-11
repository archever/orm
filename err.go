package orm

import (
	"errors"
	"log"
	"os"
)

// errors
var (
	ErrNotFund = errors.New("资源未找到")
	ErrNull    = errors.New("NUll")
)

// logger
var (
	Log   = log.New(os.Stdout, "[SQL] ", log.LstdFlags)
	Debug = log.New(os.Stdout, "[DEBUG] ", log.Lshortfile)
)
