package cli

import (
	"os"
	"sort"
	"time"

	// "github.com/urfave/cli"
)

const (
	NAME    = "app"
	VERSION = "0.0.1-beta"
	USAGE   = "Boompack service application"
)

func Run() {
	// Phase 1. Get cli options, some validation checks and configure working env
	// errors from this phase must be paniced with traceback and os.exit(1)
	// app := cli.NewApp()
	// app.Name = NAME
	// app.Version = VERSION
	// app.Compiled = time.Now()
	// app.Authors = []cli.Author{
	// 	{
	// 		Name:  "Alexander Gurinov",
	// 		Email: "alexander.gurinov@gmail.com",
	// 	},
	// 	{
	// 		Name:  "Alexey Yollov",
	// 		Email: "yollov@me.com",
	// 	},
	// }
	// app.Usage = USAGE
	// app.Flags = []cli.Flag{
	// 	cli.BoolFlag{
	// 		Name:  "debug",
	// 		Usage: debugFlagUsage,
	// 	},
	// }
	// app.Commands = []cli.Command{
	// 	{
	// 		Name:   "run",
	// 		Usage:  runCommandUsage,
	// 		Action: runCommandAction,
	// 	},
	// }
	// // configure sorting for help
	// sort.Sort(cli.FlagsByName(app.Flags))
	// sort.Sort(cli.CommandsByName(app.Commands))
	// // run
	// app.Run(os.Args)
}
