package agent

import (
	"time"

	"github.com/pathfinder-cm/pathfinder-agent/daemon"
	"github.com/pathfinder-cm/pathfinder-agent/metrics"
	"github.com/pathfinder-cm/pathfinder-agent/util"
	"github.com/pathfinder-cm/pathfinder-go-client/pfclient"
	"github.com/pathfinder-cm/pathfinder-go-client/pfmodel"
	log "github.com/sirupsen/logrus"
)

type Agent interface {
	Run(agentType string)
	RunMetrics() error
	Process() bool
	provisionContainer(sc pfmodel.Container, lcs *pfmodel.ContainerList) (bool, error)
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

func (a *agent) Run(agentType string) {
	if agentType == "metrics" {
		log.WithFields(log.Fields{}).Warn("Push Metrics")

		for {
			delay := 6 + util.RandomIntRange(1, 6)
			time.Sleep(time.Duration(delay) * time.Second)

			a.RunMetrics()
		}
		return
	}

	for {
		// Add delay between processing
		delay := 5 + util.RandomIntRange(1, 5)
		time.Sleep(time.Duration(delay) * time.Second)

		a.Process()
	}
}

func (a *agent) RunMetrics() error {
	m := metrics.NewMetrics()
	collectedMetrics := m.Collect()
	err := a.pfclient.PushMetrics(collectedMetrics)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

func (a *agent) Process() bool {
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

func (a *agent) provisionContainer(sc pfmodel.Container, lcs *pfmodel.ContainerList) (bool, error) {
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
			}).Error("Error during container creation")
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

func (a *agent) deleteContainer(sc pfmodel.Container, lcs *pfmodel.ContainerList) (bool, error) {
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
