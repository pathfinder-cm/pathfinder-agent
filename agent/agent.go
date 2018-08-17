package agent

import (
	"io/ioutil"
	"time"

	"github.com/giosakti/pathfinder-agent/daemon"
	"github.com/giosakti/pathfinder-agent/pfclient"
	log "github.com/sirupsen/logrus"
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
		serverContainers, _ := pfclient.NewContainerListFromByte(b)

		// Get from local daemon
		localContainers, err := a.containerDaemon.ListContainers()
		if err != nil {
			log.Error(err.Error())
		} else {
			// Compare containers from server and local daemon
			for _, sc := range *serverContainers {
				i := localContainers.FindByName(sc.Name)
				if i == -1 {
					a.containerDaemon.CreateContainer(sc.Name, sc.Image)
					log.WithFields(log.Fields{
						"name":   sc.Name,
						"number": sc.Image,
					}).Info("Container created")
				} else {
					localContainers.DeleteAt(i)
				}
			}
		}

		time.Sleep(5 * time.Second)
	}
}
