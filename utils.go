package orm

import (
	"fmt"
	"log"
	"os"
)

func FieldWrapper(field string) string {
	return fmt.Sprintf("%s%s%s", "`", field, "`")
}

// logger
var (
	Log   = log.New(os.Stdout, "[SQL] ", log.LstdFlags)
	Debug = log.New(os.Stdout, "[DEBUG] ", log.Lshortfile)
)
