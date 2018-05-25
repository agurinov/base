package cli

import (
	"fmt"

	"github.com/urfave/cli"
	"github.com/agurinov/dnskek/log"
)

var (
	// Actions
	runCommandUsage     = "Run application server"
	// Flags
	debugFlagUsage = "Debugging mode"
)

func runCommandAction(c *cli.Context) {
	log.SetDebug(c.GlobalBool("debug"))
	
	fmt.Println("FFFF")
}
