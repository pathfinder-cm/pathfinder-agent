package daemon

import "github.com/pathfinder-cm/pathfinder-go-client/pfmodel"

type ContainerDaemon interface {
	ListContainers() (*pfmodel.ContainerList, error)
	CreateContainer(container pfmodel.Container) (bool, string, error)
	DeleteContainer(hostname string) (bool, error)
	CreateContainerBootstrapScript(container pfmodel.Container) (bool, error)
	BootstrapContainer(container pfmodel.Container) (bool, error)
}
