package server

import (
	"path/filepath"
	"runtime"

	"github.com/boomfunc/base/server/flow"
	"github.com/boomfunc/log"
)

// TODO ADD https://stackoverflow.com/questions/17817204/how-to-set-ulimit-n-from-a-golang-program
func StartupLog(transportName, applicationName, addr, config string) {
	// calculate current working directory
	var cwd string
	cwd, err := filepath.Abs("./")
	if err != nil {
		cwd = "<unknown>"
	}

	log.Infof(
		"server:\t%s (%s application) up and running on %s",
		log.Wrap(transportName, log.Bold),
		log.Wrap(applicationName, log.Bold),
		log.Wrap(addr, log.Bold, log.Blink),
	)
	log.Infof("server:\tCurrent working directory: %s", log.Wrap(cwd, log.Bold))
	log.Infof("router:\t%s", log.Wrap(config, log.Bold))
	log.Debugf("server:\tEnabled %s mode", log.Wrap("DEBUG", log.Bold, log.Blink))
}

func PerformanceLog(numWorkers int) {
	// TODO https://insights.sei.cmu.edu/sei_blog/2017/08/multicore-and-virtualization-an-introduction.html
	log.Debugf("server:\tSpawned %d initial goroutines", runtime.NumGoroutine())
	if runtime.NumGoroutine() != numWorkers+3 {
		log.Warnf(
			"server:\tUnexpected number of initial goroutines, possibly an issue. Expected: %d, Got: %d",
			numWorkers+3,
			runtime.NumGoroutine(),
		)
	}
	log.Debugf("server:\tDetected %d CPU cores", runtime.NumCPU())
	if runtime.NumCPU() < numWorkers {
		log.Warnf(
			"server:\tPossible overloading of CPU cores. Detected: %[1]d CPU. Recommended worker number: %[1]d (Current: %[2]d)",
			runtime.NumCPU(), numWorkers,
		)
	} else if runtime.NumCPU() > numWorkers {
		log.Warnf(
			"server:\tPossible performance improvements. Increase worker number. Detected: %[1]d CPU. Recommended worker number: %[1]d (Current: %[2]d)",
			runtime.NumCPU(), numWorkers,
		)
	}
}

func AccessLog(flow *flow.Data) {
	var status, url string

	if flow.Stat.Successful() {
		status = "SUCCESS"
	} else {
		status = "ERROR"
	}

	// Request might be nil if err while parsing incoming message
	if flow.Stat.Request != nil {
		url = flow.Stat.Request.Url.RequestURI()
	} else {
		url = "/XXX/XXX/XXX"
	}

	log.Infof("%s\t-\t%s\t-\t%s\t-\tTiming: `%s`", flow.UUID.String(), url, status, flow.Chronometer)
}
