package daemon

import "github.com/pathfinder-cm/pathfinder-go-client/pfmodel"

type ContainerDaemon interface {
	ListContainers() (*pfmodel.ContainerList, error)
	CreateContainer(hostname string, image string) (bool, string, error)
	DeleteContainer(hostname string) (bool, error)
}
