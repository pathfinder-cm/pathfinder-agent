package agent

import (
	"fmt"
	"time"

	"github.com/pathfinder-cm/pathfinder-agent/daemon"
	"github.com/pathfinder-cm/pathfinder-agent/util"
	"github.com/pathfinder-cm/pathfinder-go-client/pfclient"
	"github.com/pathfinder-cm/pathfinder-go-client/pfmodel"
	log "github.com/sirupsen/logrus"
)

type bootstrapAgent struct {
	nodeHostname    string
	fullPath        string
	containerDaemon daemon.ContainerDaemon
	pfclient        pfclient.Pfclient
}

func NewBootstrapAgent(
	nodeHostname string,
	fullPath string,
	containerDaemon daemon.ContainerDaemon,
	pfclient pfclient.Pfclient) Agent {

	return &bootstrapAgent{
		nodeHostname:    nodeHostname,
		fullPath:        fullPath,
		containerDaemon: containerDaemon,
		pfclient:        pfclient,
	}
}

func (a *bootstrapAgent) Run() {
	log.WithFields(log.Fields{}).Warn("Starting bootstrap agent...")

	for {
		// Add delay between processing
		delay := 5 + util.RandomIntRange(1, 5)
		time.Sleep(time.Duration(delay) * time.Second)

		a.Process()
	}
}

func (a *bootstrapAgent) Process() bool {
	pcs, err := a.pfclient.FetchContainersFromServer(a.nodeHostname, "ListProvisionedContainers")
	if err != nil {
		return false
	}

	// Compare containers between server and local daemon
	// Do action as necessary
	for _, pc := range *pcs {
		err := a.createContainerFile(pc)
		if err != nil {
			return false
		}

		a.bootstrapContainer(pc)
	}

	return true
}

func (a *bootstrapAgent) createContainerFile(pc pfmodel.Container) error {
	log.WithFields(log.Fields{
		"hostname":      pc.Hostname,
		"ipaddress":     pc.Ipaddress,
		"source_type":   pc.Source.Type,
		"alias":         pc.Source.Alias,
		"certificate":   pc.Source.Remote.Certificate,
		"mode":          pc.Source.Mode,
		"server":        pc.Source.Remote.Server,
		"protocol":      pc.Source.Remote.Protocol,
		"auth_type":     pc.Source.Remote.AuthType,
		"bootstrappers": pc.Bootstrappers,
	}).Info("Creating container file")

	err := a.containerDaemon.CreateContainerFile(pc, a.fullPath)
	if err != nil {
		a.pfclient.MarkContainerAsBootstrapError(
			a.nodeHostname,
			pc.Hostname,
		)
		log.WithFields(log.Fields{
			"hostname":      pc.Hostname,
			"ipaddress":     pc.Ipaddress,
			"source_type":   pc.Source.Type,
			"alias":         pc.Source.Alias,
			"certificate":   pc.Source.Remote.Certificate,
			"mode":          pc.Source.Mode,
			"server":        pc.Source.Remote.Server,
			"protocol":      pc.Source.Remote.Protocol,
			"auth_type":     pc.Source.Remote.AuthType,
			"bootstrappers": pc.Bootstrappers,
		}).Error(fmt.Sprintf("Error during creating container file. %v", err))
		return err
	}

	return nil
}

func (a *bootstrapAgent) bootstrapContainer(pc pfmodel.Container) (bool, error) {
	log.WithFields(log.Fields{
		"hostname":      pc.Hostname,
		"ipaddress":     pc.Ipaddress,
		"source_type":   pc.Source.Type,
		"alias":         pc.Source.Alias,
		"certificate":   pc.Source.Remote.Certificate,
		"mode":          pc.Source.Mode,
		"server":        pc.Source.Remote.Server,
		"protocol":      pc.Source.Remote.Protocol,
		"auth_type":     pc.Source.Remote.AuthType,
		"bootstrappers": pc.Bootstrappers,
	}).Info("Bootstrapping container")

	ok, err := a.containerDaemon.ExecContainer(pc, a.fullPath)
	if !ok {
		a.pfclient.MarkContainerAsBootstrapError(
			a.nodeHostname,
			pc.Hostname,
		)
		log.WithFields(log.Fields{
			"hostname":      pc.Hostname,
			"ipaddress":     pc.Ipaddress,
			"source_type":   pc.Source.Type,
			"alias":         pc.Source.Alias,
			"certificate":   pc.Source.Remote.Certificate,
			"mode":          pc.Source.Mode,
			"server":        pc.Source.Remote.Server,
			"protocol":      pc.Source.Remote.Protocol,
			"auth_type":     pc.Source.Remote.AuthType,
			"bootstrappers": pc.Bootstrappers,
		}).Error(fmt.Sprintf("Error during container bootsrapping. %v", err))
		return false, err
	}

	a.pfclient.MarkContainerAsBootstrapped(
		a.nodeHostname,
		pc.Hostname,
	)
	log.WithFields(log.Fields{
		"hostname":                pc.Hostname,
		"ipaddress":               pc.Ipaddress,
		"source_type":             pc.Source.Type,
		"alias":                   pc.Source.Alias,
		"certificate":             pc.Source.Remote.Certificate,
		"mode":                    pc.Source.Mode,
		"server":                  pc.Source.Remote.Server,
		"protocol":                pc.Source.Remote.Protocol,
		"auth_type":               pc.Source.Remote.AuthType,
		"bootstrap_type":          pc.Bootstrappers[0].Type,
		"bootstrap_cookbooks_url": pc.Bootstrappers[0].CookbooksUrl,
		"bootstrap_attributes":    pc.Bootstrappers[0].Attributes,
	}).Info("Container bootstrapped")

	return true, nil
}
