package main

import (
	"net"
	"net/http"
	"os"
	"time"

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
	hostname, _ := os.Hostname()
	daemon := daemon.LXD{SocketPath: DefaultLXDSocketPath}
	client := &http.Client{
		Timeout: time.Second * 60,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 60 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 60 * time.Second,
		},
	}
	a := agent.NewAgent(
		hostname,
		daemon,
		client,
		DefaultPfServerAddr,
		DefaultListContainersPath,
		DefaultProvisionedPath,
	)
	a.Run()
}
