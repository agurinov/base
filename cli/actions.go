package cli

import (
	"net"
	"os"

	"github.com/boomfunc/log"
	"github.com/urfave/cli"

	"github.com/boomfunc/base/server"
)

var (
	// Actions
	runCommandUsage     = "Run creates concrete type of server and listen for incoming requests. Choose subcommand"
	runUDPCommandUsage  = "Run UDP application server"
	runTCPCommandUsage  = "Run TCP application server"
	runHTTPCommandUsage = "Run HTTP application server"
	// Flags
	debugFlagUsage  = "Debugging mode"
	portFlagUsage   = "Port on which the listener will be"
	configFlagUsage = "Path to config file"
)

func runTCPCommandAction(c *cli.Context) {
	log.SetDebug(c.GlobalBool("debug"))

	server, err := server.NewTCP(
		net.ParseIP("0.0.0.0"),
		c.GlobalInt("port"),
		c.GlobalString("config"),
	)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	server.Serve()
}
