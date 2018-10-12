package dispatcher

import (
	"github.com/boomfunc/log"
)

func StartupLog(workerNum int) {
	log.Debugf("dispatcher:\tSpawned %s worker(s)", log.Wrap(workerNum, log.Bold))
}
