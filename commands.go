package main

import (
	"net"
	"net/http"
	"os"
	"time"

	"github.com/BaritoLog/go-boilerplate/srvkit"
	"github.com/pathfinder-cm/pathfinder-agent/agent"
	"github.com/pathfinder-cm/pathfinder-agent/daemon"
	"github.com/pathfinder-cm/pathfinder-go-client/pfclient"
	"github.com/pathfinder-cm/pathfinder-agent/util"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

//"errors"
	

func CmdAgent(ctx *cli.Context) {
	if ctx.Bool("verbose") == true {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	log.WithFields(log.Fields{
		"verbose": ctx.Bool("verbose"),
	}).Warn("Agent starting...")
	runAgent()
	srvkit.GracefullShutdown(func() {
		log.Warn("Agent stopping...")
	})
}


func runAgent() {
	hostname, _ := os.Hostname()
	ipaddress := getLocalIP()
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

	for {
		log.WithFields(log.Fields{}).Warn("Trying to register to pathfinder server...")
		ok, _ := pfclient.Register(hostname,ipaddress)
		
		if !ok {
			log.Error("Cannot register to pathfinder server, please check your configuration")

			delay := 60 + util.RandomIntRange(1, 10)
			time.Sleep(time.Duration(delay) * time.Second)			
		} else {
			break
		}
	}

	provisionAgent := agent.NewProvisionAgent(hostname, daemon, pfclient)
	go provisionAgent.Run()

	metricsAgent := agent.NewMetricsAgent(hostname, pfclient)
	go metricsAgent.Run()
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
