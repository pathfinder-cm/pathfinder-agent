package main

import (
	"github.com/giosakti/pathfinder-agent/agent"
	"github.com/urfave/cli"
)

func CmdAgent(ctx *cli.Context) {
	go runAgent()
}

func runAgent() {
	a := agent.NewAgent()
	a.Run()
}
