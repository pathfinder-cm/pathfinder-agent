package main

import (
	"fmt"

	"github.com/BaritoLog/go-boilerplate/srvkit"
	"github.com/giosakti/pathfinder-agent/agent"
	"github.com/urfave/cli"
)

func CmdAgent(ctx *cli.Context) {
	fmt.Println("Agent starting...")
	go runAgent()
	srvkit.GracefullShutdown(func() {
		fmt.Println("Agent stopping...")
	})
}

func runAgent() {
	a := agent.NewAgent()
	a.Run()
}
