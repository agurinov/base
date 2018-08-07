package dispatcher

import (
	"github.com/boomfunc/log"
)

func StartupLog(workerNum int) {
	log.Infof("dispatcher:\tSpawned %s workers", log.Wrap(workerNum, log.Bold))
}
