package main

import (
	"github.com/BaritoLog/go-boilerplate/srvkit"
	"github.com/giosakti/pathfinder-agent/agent"
	"github.com/giosakti/pathfinder-agent/daemon"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func CmdAgent(ctx *cli.Context) {
	if ctx.Bool("verbose") == true {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	log.WithFields(log.Fields{
		"verbose": ctx.Bool("verbose"),
	}).Warn("Agent starting...")
	go runAgent()
	srvkit.GracefullShutdown(func() {
		log.Warn("Agent stopping...")
	})
}

func runAgent() {
	d := daemon.LXD{SocketPath: DefaultLXDSocketPath}
	a := agent.NewAgent(d)
	a.Run()
}
