package cli

import (
	"fmt"
	"net"
	"os"
	"strings"

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

func runCommandAction(c *cli.Context) {
	log.SetDebug(c.GlobalBool("debug"))

	StartupLog(c.App.Version, c.App.Compiled)

	// Exctract params
	transport := c.Command.Name
	ip := net.ParseIP("0.0.0.0")
	port := c.GlobalInt("port")
	filename := c.GlobalString("config")

	// Create server
	srv, err := server.New(transport, ip, port, filename)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// Run
	server.StartupLog(strings.ToUpper(transport), fmt.Sprintf("%s:%d", ip, port), filename)
	srv.Serve()
}
