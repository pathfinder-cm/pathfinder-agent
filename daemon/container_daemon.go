package daemon

import "github.com/pathfinder-cm/pathfinder-go-client/pfmodel"

type ContainerDaemon interface {
	ListContainers() (*pfmodel.ContainerList, error)
	CreateContainer(hostname string, image_alias string, image_server string, image_protocol string) (bool, string, error)
	DeleteContainer(hostname string) (bool, error)
}
