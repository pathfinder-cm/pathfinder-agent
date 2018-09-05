package agent

import (
	"fmt"
	"github.com/pathfinder-cm/pathfinder-agent/daemon"
	"github.com/pathfinder-cm/pathfinder-agent/util"
	"github.com/pathfinder-cm/pathfinder-go-client/pfclient"
	"github.com/pathfinder-cm/pathfinder-go-client/pfmodel"
	log "github.com/sirupsen/logrus"
	"time"
)

type provisionAgent struct {
	nodeHostname    string
	containerDaemon daemon.ContainerDaemon
	pfclient        pfclient.Pfclient
}

func (a *provisionAgent) Run() {
	for {
		// Add delay between processing
		delay := 5 + util.RandomIntRange(1, 5)
		time.Sleep(time.Duration(delay) * time.Second)

		a.Process()
	}
}

func (a *provisionAgent) Process() bool {
	scs, err := a.pfclient.FetchContainersFromServer(a.nodeHostname)
	if err != nil {
		return false
	}

	// Compare containers between server and local daemon
	// Do action as necessary
	for _, sc := range *scs {
		lcs, err := a.containerDaemon.ListContainers()
		if err != nil {
			return false
		}

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

func (a *provisionAgent) provisionContainer(sc pfmodel.Container, lcs *pfmodel.ContainerList) (bool, error) {
	i := lcs.FindByHostname(sc.Hostname)
	if i == -1 {
		log.WithFields(log.Fields{
			"hostname": sc.Hostname,
			"image":    sc.Image,
		}).Info("Creating container")

		ok, ipaddress, err := a.containerDaemon.CreateContainer(sc.Hostname, sc.Image)
		if !ok {
			a.pfclient.MarkContainerAsProvisionError(
				a.nodeHostname,
				sc.Hostname,
			)
			log.WithFields(log.Fields{
				"hostname": sc.Hostname,
				"image":    sc.Image,
			}).Error(fmt.Sprintf("Error during container creation. %v", err))
			return false, err
		}

		a.pfclient.UpdateIpaddress(
			a.nodeHostname,
			sc.Hostname,
			ipaddress,
		)
		a.pfclient.MarkContainerAsProvisioned(
			a.nodeHostname,
			sc.Hostname,
		)

		log.WithFields(log.Fields{
			"hostname": sc.Hostname,
			"image":    sc.Image,
		}).Info("Container created")
	} else {
		log.WithFields(log.Fields{
			"hostname": sc.Hostname,
			"image":    sc.Image,
		}).Info("Container already exist")

		a.pfclient.MarkContainerAsProvisioned(
			a.nodeHostname,
			sc.Hostname,
		)
	}

	return true, nil
}

func (a *provisionAgent) deleteContainer(sc pfmodel.Container, lcs *pfmodel.ContainerList) (bool, error) {
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
