package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	Name                      = "Pathfinder Agent"
	Version                   = "0.0.1"
	PfCluster                 = "default"
	DefaultLXDSocketPath      = "/var/snap/lxd/common/lxd/unix.socket"
	DefaultPfServerAddr       = "http://192.168.33.1:3000"
	DefaultListContainersPath = "api/v1/node/containers/scheduled"
	DefaultProvisionedPath    = "api/v1/node/containers/provision"
	DefaultProvisionErrorPath = "api/v1/node/containers/provision_error"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
}

func main() {
	app := cli.App{
		Name:    Name,
		Usage:   "Agent for Pathfinder container manager",
		Version: Version,
		Action:  CmdAgent,
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose, V",
			Usage: "Enable verbose mode",
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err.Error())
	}
}
