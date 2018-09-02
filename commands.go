package main

import (
	"net"
	"net/http"
	"os"
	"time"

	"github.com/BaritoLog/go-boilerplate/srvkit"
	"github.com/pathfinder-cm/pathfinder-agent/agent"
	"github.com/pathfinder-cm/pathfinder-agent/daemon"
	"github.com/pathfinder-cm/pathfinder-agent/pfclient"
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
	daemon, err := daemon.NewLXD(hostname, LXDSocketPath)
	if err != nil {
		log.Error("Cannot connect to container daemon")
		return
	}
	httpClient := &http.Client{
		Timeout: time.Second * 60,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 60 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 60 * time.Second,
		},
	}
	pfclient := pfclient.NewPfclient(
		PfCluster,
		PfClusterPassword,
		httpClient,
		PfServerAddr,
		PfApiPath,
	)

	a := agent.NewAgent(hostname, daemon, pfclient)
	a.Run()
}
