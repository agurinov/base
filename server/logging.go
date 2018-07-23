package server

import (
	"runtime/debug"

	"github.com/boomfunc/base/server/request"
	"github.com/boomfunc/log"
)

func StartupLog(mode, addr, filename string) {
	log.Infof("%s server up and running on %s", log.Wrap(mode, log.Bold), log.Wrap(addr, log.Bold, log.Blink))
	log.Infof("Spawned config file: %s", log.Wrap(filename, log.Bold))
	log.Debugf("Enabled %s mode", log.Wrap("DEBUG", log.Bold, log.Blink))
}

func AccessLog(stat request.Stat) {
	req := stat.Request
	status := "SUCCESS"
	if !stat.Successful() {
		status = "ERROR"
	}

	log.Infof("%s\t-\t%s\t-\t%s\t-\t%s", req.UUID(), req.Url(), status, stat.Duration)
}

// TODO clear Stack
func ErrorLog(err interface{}) {
	log.Errorf("%s\n%s", err, debug.Stack())
}
