package cli

import (
	"github.com/boomfunc/log"
	"github.com/urfave/cli"
)

var (
	// Actions
	runCommandUsage = "Run application server"
	// Flags
	debugFlagUsage  = "Debugging mode"
	strictFlagUsage = `Strict mode. If any of the following conditions is not satisfied there will be an error
	1. Config is invalid yaml`
)

func runCommandAction(c *cli.Context) {
	log.SetDebug(c.GlobalBool("debug"))

	log.Info("Info")
	log.Debug("Debug")
}
