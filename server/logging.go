package server

import (
	"runtime/debug"

	"github.com/boomfunc/log"
)

func serverStartupLog(mode, addr, filename string) {
	log.Infof("%s server up and running on %s", mode, log.Wrap(addr, log.Bold, log.Blink))
	log.Infof("Spawned config file: %s", log.Wrap(filename, log.Bold))
	log.Debugf("Enabled %s mode", log.Wrap("DEBUG", log.Bold, log.Blink))
}

func serverAccessLog(req Request, status string) {
	log.Infof("%s\t-\t%s\t-\t%s", req.UUID(), req.Url(), status)
}

// TODO clear Stack
func serverErrorLog(err interface{}) {
	log.Errorf("%s\n%s", err, debug.Stack())
}