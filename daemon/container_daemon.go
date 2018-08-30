package daemon

import (
	"github.com/pathfinder-cm/pathfinder-agent/model"
)

type ContainerDaemon interface {
	ListContainers() (*model.ContainerList, error)
	CreateContainer(hostname string, image string) (bool, string, error)
	DeleteContainer(hostname string) (bool, error)
}
