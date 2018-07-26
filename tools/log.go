package tools

import (
	"os"
	"runtime/debug"

	"github.com/boomfunc/log"
)

// TODO clear Stack
func ErrorLog(err interface{}) {
	log.Errorf("%s\n%s", err, debug.Stack())
}

func FatalLog(err interface{}) {
	ErrorLog(err)
	os.Exit(1)
}
