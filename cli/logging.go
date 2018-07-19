package cli

import (
	"time"

	"github.com/boomfunc/log"
)

func StartupLog(version string, compiled time.Time) {
	log.Infof("************************************************************")
	log.Infof("Boomfunc base version:\t%s", log.Wrap(version, log.Bold))
	log.Infof("Boomfunc compilation time:\t%s", log.Wrap(compiled.String(), log.Bold))
	log.Infof("************************************************************")
	log.Infof("")
}
