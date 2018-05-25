package cli

import (
	"github.com/boomfunc/log"
	"github.com/urfave/cli"
)

var (
	// Actions
	runCommandUsage = "Run application server"
	// Flags
	debugFlagUsage = "Debugging mode"
)

func runCommandAction(c *cli.Context) {
	log.SetDebug(c.GlobalBool("debug"))

	log.Info("Info")
	log.Debug("Debug")
}
