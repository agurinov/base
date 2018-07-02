package cli

import (
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/urfave/cli"
)

const (
	NAME  = "base"
	USAGE = "Boompack service application"
)

func Run(VERSION, TIMESTAMP string) {
	// prepare build variables passed through -ldflags
	ts, err := strconv.ParseInt(TIMESTAMP, 10, 64)
	if err != nil {
		panic(err)
	}
	// Phase 1. Get cli options, some validation checks and configure working env
	// errors from this phase must be paniced with traceback and os.exit(1)
	app := cli.NewApp()
	app.Name = NAME
	app.Version = VERSION
	app.Compiled = time.Unix(ts, 0)
	app.Authors = []cli.Author{
		{
			Name:  "Alexander Gurinov",
			Email: "alexander.gurinov@gmail.com",
		},
		{
			Name:  "Alexey Yollov",
			Email: "yollov@me.com",
		},
	}
	app.Usage = USAGE
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "debug",
			Usage:  debugFlagUsage,
			EnvVar: "BMP_DEBUG_MODE",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:  "run",
			Usage: runCommandUsage,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:   "port",
					Usage:  portFlagUsage,
					EnvVar: "BMP_PORT",
					Value:  8080,
				},
				cli.StringFlag{
					Name:   "config",
					Usage:  configFlagUsage,
					EnvVar: "BMP_CONFIG",
					Value:  "/boomfunc/app/conf.yml",
				},
			},
			Subcommands: []cli.Command{
				{
					Name:   "udp",
					Usage:  runUDPCommandUsage,
					Action: runTCPCommandAction,
				},
				{
					Name:   "tcp",
					Usage:  runTCPCommandUsage,
					Action: runTCPCommandAction,
				},
				{
					Name:   "http",
					Usage:  runHTTPCommandUsage,
					Action: runTCPCommandAction,
				},
			},
		},
	}
	// configure sorting for help
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	// run
	app.Run(os.Args)
}
