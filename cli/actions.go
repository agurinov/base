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
	portFlagUsage   = "Port on which the listener will be"
	configFlagUsage = "Path to config file"
)

func runCommandAction(c *cli.Context) {
	log.SetDebug(c.GlobalBool("debug"))

	server, err := server.NewTCP(
		net.ParseIP("0.0.0.0"),
		c.Int("port"),
		c.String("config"),
	)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	server.Serve()
}
