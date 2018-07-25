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
	var status, uuid, url string

	if stat.Successful() {
		status = "SUCCESS"
	} else {
		status = "ERROR"
	}

	// Request might be nil if err while parsing incoming message
	if stat.Request != nil {
		uuid = stat.Request.UUID.String()
		url = stat.Request.Url
	} else {
		uuid = "<not_parsed>"
		url = "<not_parsed>"
	}

	log.Infof("%s\t-\t%s\t-\t%s\t-\t%s\t-\tWritten: %d", uuid, url, status, stat.Duration, stat.Len)
}

// TODO clear Stack
func ErrorLog(err interface{}) {
	log.Errorf("%s\n%s", err, debug.Stack())
}
