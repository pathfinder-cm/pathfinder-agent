package agent

import (
	"fmt"
	"time"

	"github.com/giosakti/pathfinder-agent/daemon"
	"github.com/giosakti/pathfinder-agent/model"
	"github.com/giosakti/pathfinder-agent/pfclient"
	"github.com/giosakti/pathfinder-agent/util"
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
		delay := 5 + util.RandomIntRange(1, 5)
		time.Sleep(time.Duration(delay) * time.Second)

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
		switch status := sc.Status; status {
		case "SCHEDULED":
			a.provisionContainer(sc, lcs)
		case "SCHEDULE_DELETION":
			a.deleteContainer(sc, lcs)
		}
	}

	//TODO: get a list of orphaned containers

	return true
}

func (a *agent) provisionContainer(sc model.Container, lcs *model.ContainerList) (bool, error) {
	i := lcs.FindByHostname(sc.Hostname)
	if i == -1 {
		log.WithFields(log.Fields{
			"hostname": sc.Hostname,
			"image":    sc.Image,
		}).Info("Creating container")

		ok, err := a.containerDaemon.CreateContainer(sc.Hostname, sc.Image)
		if !ok {
			a.pfclient.MarkContainerAsProvisionError(
				a.nodeHostname,
				sc.Hostname,
			)
			log.WithFields(log.Fields{
				"hostname": sc.Hostname,
				"image":    sc.Image,
			}).Error("Error during container creation")
			return false, err
		}

		ok, err = a.pfclient.MarkContainerAsProvisioned(
			a.nodeHostname,
			sc.Hostname,
		)
		if !ok {
			return false, err
		}

		log.WithFields(log.Fields{
			"hostname": sc.Hostname,
			"image":    sc.Image,
		}).Info("Container created")
	} else {
		log.WithFields(log.Fields{
			"hostname": sc.Hostname,
			"image":    sc.Image,
		}).Info("Container already exist")
	}

	return true, nil
}

func (a *agent) deleteContainer(sc model.Container, lcs *model.ContainerList) (bool, error) {
	i := lcs.FindByHostname(sc.Hostname)
	if i == -1 {
		log.WithFields(log.Fields{
			"hostname": sc.Hostname,
			"image":    sc.Image,
		}).Info("Container already deleted")
	} else {
		ok, err := a.containerDaemon.DeleteContainer(sc.Hostname)
		if !ok {
			log.WithFields(log.Fields{
				"hostname": sc.Hostname,
				"image":    sc.Image,
			}).Error("Error during container deletion")
			fmt.Println(err)
			return false, err
		}

		log.WithFields(log.Fields{
			"hostname": sc.Hostname,
			"image":    sc.Image,
		}).Info("Container deleted")
	}

	a.pfclient.MarkContainerAsDeleted(
		a.nodeHostname,
		sc.Hostname,
	)

	return true, nil
}
