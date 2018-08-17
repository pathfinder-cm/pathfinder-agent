package main

import (
	"fmt"

	"github.com/BaritoLog/go-boilerplate/srvkit"
	"github.com/giosakti/pathfinder-agent/agent"
	"github.com/giosakti/pathfinder-agent/daemon"
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
	d := daemon.LXD{SocketPath: DefaultLXDSocketPath}
	a := agent.NewAgent(d)
	a.Run()
}
