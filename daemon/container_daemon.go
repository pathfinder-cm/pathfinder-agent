package daemon

import "github.com/pathfinder-cm/pathfinder-go-client/pfmodel"

type ContainerDaemon interface {
	ListContainers() (*pfmodel.ContainerList, error)
	CreateContainer(hostname string, source_type string, alias string, certificate string, mode string, server string, protocol string) (bool, string, error)
	DeleteContainer(hostname string) (bool, error)
}
