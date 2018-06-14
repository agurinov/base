package cli

import (
	"net"
	"os"

	"github.com/boomfunc/log"
	"github.com/urfave/cli"

	"app/server"
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

	server, err := server.NewTCP(
		net.ParseIP("0.0.0.0"),
		8080,
		"./conf/test.yml",
	)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	server.Serve()
}
