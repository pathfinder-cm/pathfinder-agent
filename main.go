package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

const (
	Name    = "Pathfinder Agent"
	Version = "0.0.1"
)

func main() {
	app := cli.App{
		Name:    Name,
		Usage:   "Agent for Pathfinder container manager",
		Version: Version,
		Action:  CmdAgent,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(fmt.Sprintf("Fatal error: %s", err.Error()))
	}
}
