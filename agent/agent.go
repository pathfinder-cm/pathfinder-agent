package agent

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/giosakti/pathfinder-agent/daemon"
	"github.com/giosakti/pathfinder-agent/pfclient"
)

type Agent interface {
	Run()
}

type agent struct {
	containerDaemon daemon.ContainerDaemon
}

func NewAgent(cd daemon.ContainerDaemon) Agent {
	return &agent{
		containerDaemon: cd,
	}
}

func (a *agent) Run() {
	for {
		// Get from API Server
		b, _ := ioutil.ReadFile("/opt/projects/golang/src/github.com/giosakti/pathfinder-agent/fixtures/scheduled-containers.json")
		remoteContainers, _ := pfclient.NewContainerListFromByte(b)

		// Get from local daemon
		localContainers, _ := a.containerDaemon.ListContainers()

		// Compare containers from server and local daemon
		for _, rc := range *remoteContainers {
			i := localContainers.FindByName(rc.Name)
			if i == -1 {
				fmt.Println("Creating Container", rc.Name)
				a.containerDaemon.CreateContainer(rc.Name, rc.Image)
			} else {
				localContainers.DeleteAt(i)
			}
		}

		time.Sleep(5 * time.Second)
	}
}
