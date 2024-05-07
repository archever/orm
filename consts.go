package orm

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

var (
	_, b, _, _ = runtime.Caller(0)
	BaseDir    = filepath.Clean(fmt.Sprintf("%s", filepath.Dir(b)))
)

func init() {
	err := godotenv.Load(fmt.Sprintf("%s/.env", BaseDir))
	if err != nil {
		Log.Printf("load .env failed")
	}
}
