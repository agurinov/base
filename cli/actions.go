package cli

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"text/template"

	"github.com/urfave/cli"

	"github.com/agurinov/dnskek/log"
	"github.com/agurinov/dnskek/server"
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
