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

func AccessLog(response request.Response) {
	req := response.Request

	log.Infof("%s\t-\t%s\t-\t%s", req.UUID(), req.Url(), response.Status)
}

// TODO clear Stack
func ErrorLog(err interface{}) {
	log.Errorf("%s\n%s", err, debug.Stack())
}
