package daemon

import (
	"github.com/giosakti/pathfinder-agent/model"
)

type ContainerDaemon interface {
	ListContainers() (*model.ContainerList, error)
	CreateContainer(name string, image string) (bool, error)
}
