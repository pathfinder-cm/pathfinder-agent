package daemon

import (
	"github.com/giosakti/pathfinder-agent/model"
)

type ContainerDaemon interface {
	ListContainers() (*model.ContainerList, error)
	CreateContainer(hostname string, image string) (bool, error)
	DeleteContainer(hostname string) (bool, error)
}
