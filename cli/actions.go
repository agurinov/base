package cli

import (
	"fmt"
	"net"
	"strings"

	"github.com/boomfunc/base/server"
	"github.com/boomfunc/base/tools"
	"github.com/boomfunc/log"
	"github.com/urfave/cli"
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
	appFlagUsage    = "Variant of server application layer"
)

func runCommandAction(c *cli.Context) {
	log.SetDebug(c.GlobalBool("debug"))

	StartupLog("base", c.App.Version, c.App.Compiled)

	// Exctract params
	transport := c.Command.Name
	application := c.GlobalString("app")
	ip := net.ParseIP("0.0.0.0")
	port := c.GlobalInt("port")
	workerNum := c.GlobalInt("workers")
	filename := c.GlobalString("config")

	// Create server
	srv, err := server.New(transport, application, ip, port, filename)
	if err != nil {
		tools.FatalLog(err)
	}

	// Run
	server.StartupLog(strings.ToUpper(transport), strings.ToUpper(application), fmt.Sprintf("%s:%d", ip, port), filename)
	srv.Serve(workerNum)
}
