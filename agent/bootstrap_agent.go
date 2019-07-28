package agent

import (
	"fmt"

	"github.com/pathfinder-cm/pathfinder-agent/daemon"
	"github.com/pathfinder-cm/pathfinder-go-client/pfclient"
	"github.com/pathfinder-cm/pathfinder-go-client/pfmodel"
	log "github.com/sirupsen/logrus"
)

type bootstrapAgent struct {
	nodeHostname    string
	containerDaemon daemon.ContainerDaemon
	pfclient        pfclient.Pfclient
}

func NewBootstrapAgent(
	nodeHostname string,
	containerDaemon daemon.ContainerDaemon,
	pfclient pfclient.Pfclient) Agent {

	return &bootstrapAgent{
		nodeHostname:    nodeHostname,
		containerDaemon: containerDaemon,
		pfclient:        pfclient,
	}
}

func (a *bootstrapAgent) Run() {

}

func (a *bootstrapAgent) Process() bool {
	pcs, err := a.pfclient.FetchProvisionedContainersFromServer(a.nodeHostname)
	if err != nil {
		return false
	}

	for _, pc := range *pcs {
		err := a.createContainerBootstrapScript(pc)
		if err != nil {
			return false
		}

		a.bootstrapContainer(pc)
	}

	return true
}

func (a *bootstrapAgent) createContainerBootstrapScript(pc pfmodel.Container) error {
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
	}).Info("Creating container bootstrap script")

	err := a.containerDaemon.CreateContainerBootstrapScript(pc)
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
		}).Error(fmt.Sprintf("Error when creating container bootstrap script: %v", err))
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

	ok, err := a.containerDaemon.ExecContainerBootstrap(pc)
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
		}).Error(fmt.Sprintf("Error when bootstrapping container: %v", err))
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
