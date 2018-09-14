package server

import (
	"path/filepath"
	"runtime"

	"github.com/boomfunc/base/server/flow"
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

func PerformanceLog(numWorkers int) {
	// TODO https://insights.sei.cmu.edu/sei_blog/2017/08/multicore-and-virtualization-an-introduction.html
	log.Debugf("Spawned %d initial goroutines", runtime.NumGoroutine())
	if runtime.NumGoroutine() != numWorkers+2 {
		log.Warnf(
			"Unexpected number of initial goroutines, possibly an issue. Expected: %d, Got: %d",
			numWorkers+2,
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

func AccessLog(flow *flow.Data) {
	var status, uuid, url string

	if flow.Stat.Successful() {
		status = "SUCCESS"
	} else {
		status = "ERROR"
	}

	// Request might be nil if err while parsing incoming message
	if flow.Stat.Request != nil {
		uuid = flow.Stat.Request.UUID.String()
		url = flow.Stat.Request.Url.RequestURI()
	} else {
		uuid = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
		url = "/XXX/XXX/XXX"
	}

	log.Infof("%s\t-\t%s\t-\t%s\t-\tTiming: `%s`\t-\tWritten: %d", uuid, url, status, flow.Timing, flow.Stat.Len)
}
