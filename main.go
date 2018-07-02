package main

import (
	"github.com/boomfunc/base/cli"
)

var VERSION string
var TIMESTAMP string

func main() {
	cli.Run(VERSION, TIMESTAMP)
}
