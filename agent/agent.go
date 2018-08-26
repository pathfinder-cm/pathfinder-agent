package agent

import (
	"time"

	"github.com/giosakti/pathfinder-agent/daemon"
	"github.com/giosakti/pathfinder-agent/model"
	"github.com/giosakti/pathfinder-agent/pfclient"
	log "github.com/sirupsen/logrus"
)

type Agent interface {
	Run()
	Process() bool
	provisionContainer(sc model.Container, lcs *model.ContainerList) (bool, error)
}

type agent struct {
	nodeHostname    string
	containerDaemon daemon.ContainerDaemon
	pfclient        pfclient.Pfclient
}

func NewAgent(
	nodeHostname string,
	containerDaemon daemon.ContainerDaemon,
	pfclient pfclient.Pfclient) Agent {

	return &agent{
		nodeHostname:    nodeHostname,
		containerDaemon: containerDaemon,
		pfclient:        pfclient,
	}
}

func (a *agent) Run() {
	for {
		// Add delay between processing
		time.Sleep(5 * time.Second)

		a.Process()
	}
}

func (a *agent) Process() bool {
	scs, err := a.pfclient.FetchContainersFromServer(a.nodeHostname)
	if err != nil {
		return false
	}

	lcs, err := a.containerDaemon.ListContainers()
	if err != nil {
		return false
	}

	// Compare containers between server and local daemon
	// Do action as necessary
	for _, sc := range *scs {
		ok, _ := a.provisionContainer(sc, lcs)
		if !ok {
			return false
		}
	}

	return true
}

func (a *agent) provisionContainer(sc model.Container, lcs *model.ContainerList) (bool, error) {
	i := lcs.FindByHostname(sc.Hostname)
	if i == -1 {
		log.WithFields(log.Fields{
			"hostname": sc.Hostname,
			"number":   sc.Image,
		}).Info("Creating container")
		a.containerDaemon.CreateContainer(sc.Hostname, sc.Image)

		ok, err := a.pfclient.MarkContainerAsProvisioned(
			a.nodeHostname,
			sc.Hostname,
		)
		if !ok {
			return false, err
		}

		log.WithFields(log.Fields{
			"hostname": sc.Hostname,
			"number":   sc.Image,
		}).Info("Container created")
	} else {
		lcs.DeleteAt(i)
		log.WithFields(log.Fields{
			"hostname": sc.Hostname,
			"number":   sc.Image,
		}).Info("Container already exist")
	}

	return true, nil
}
