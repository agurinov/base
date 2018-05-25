package cli

import (
	"fmt"

	"github.com/agurinov/dnskek/log"
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

	fmt.Println("FFFF")
}
