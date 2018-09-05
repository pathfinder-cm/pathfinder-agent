package agent

import (
	"github.com/pathfinder-cm/pathfinder-agent/daemon"
	"github.com/pathfinder-cm/pathfinder-go-client/pfclient"
)

type Agent interface {
	Run()
}

func NewAgent(
	nodeHostname string,
	containerDaemon daemon.ContainerDaemon,
	pfclient pfclient.Pfclient,
	agentType string) Agent {

	if agentType == "provision" {
		return &provisionAgent{
			nodeHostname:    nodeHostname,
			containerDaemon: containerDaemon,
			pfclient:        pfclient,
		}
	} else {
		return &metricsAgent{
			nodeHostname:    nodeHostname,
			containerDaemon: containerDaemon,
			pfclient:        pfclient,
		}
	}
}
