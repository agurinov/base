package server

import (
	"path/filepath"
	"runtime"
	"runtime/debug"

	"github.com/boomfunc/base/server/request"
	"github.com/boomfunc/log"
)

func StartupLog(transportName, applicationName, addr, filename string) {
	var fpath string

	fpath, err := filepath.Abs(filename)
	if err != nil {
		fpath = filename
	}

	log.Infof(
		"%s server (%s application) up and running on %s",
		log.Wrap(transportName, log.Bold),
		log.Wrap(applicationName, log.Bold),
		log.Wrap(addr, log.Bold, log.Blink),
	)
	log.Infof("Spawned config file: %s", log.Wrap(fpath, log.Bold))
	log.Debugf("Enabled %s mode", log.Wrap("DEBUG", log.Bold, log.Blink))
}

func PerfomanceLog(numWorkers int) {
	// TODO https://insights.sei.cmu.edu/sei_blog/2017/08/multicore-and-virtualization-an-introduction.html
	log.Debugf("Spawned %d goroutines", runtime.NumGoroutine())
	if runtime.NumGoroutine() != numWorkers+3 {
		log.Warnf(
			"Unexpected number of goroutines, possibly an issue. Expected: %d, Got: %d",
			numWorkers+3,
			runtime.NumGoroutine(),
		)
	}
	log.Debugf("Detected %d CPU cores", runtime.NumCPU())
	if runtime.NumCPU() < numWorkers {
		log.Warnf(
			"Possible overloading of CPU cores. Detected: %[1]d CPU. Recommended worker number: %[1]d (Current: %[2]d)",
			runtime.NumCPU(), numWorkers,
		)
	} else if runtime.NumCPU() > numWorkers {
		log.Warnf(
			"Possible performance improvements. Increase worker number. Detected: %[1]d CPU. Recommended worker number: %[1]d (Current: %[2]d)",
			runtime.NumCPU(), numWorkers,
		)
	}
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
		uuid = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
		url = "/XXX/XXX/XXX"
	}

	log.Infof("%s\t-\t%s\t-\t%s\t-\t%s\t-\tWritten: %d", uuid, url, status, stat.Duration, stat.Len)
}

// TODO clear Stack
func ErrorLog(err interface{}) {
	log.Errorf("%s\n%s", err, debug.Stack())
}
