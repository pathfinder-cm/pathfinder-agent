package daemon

import (
	"github.com/giosakti/pathfinder-agent/model"
	client "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
)

type apiContainers []api.Container

type LXD struct {
	conn client.ContainerServer
}

func (a apiContainers) ToContainerList() *model.ContainerList {
	containerList := make(model.ContainerList, len(a))
	for i, c := range a {
		containerList[i] = model.Container{
			Hostname: c.Name,
		}
	}
	return &containerList
}

func NewLXD(socketPath string) (*LXD, error) {
	conn, err := client.ConnectLXDUnix(socketPath, nil)
	if err != nil {
		return nil, err
	}
	return &LXD{conn: conn}, nil
}

func (l *LXD) ListContainers() (*model.ContainerList, error) {
	res, err := l.conn.GetContainers()
	if err != nil {
		return nil, err
	}

	containerList := apiContainers(res).ToContainerList()

	return containerList, nil
}

func (l *LXD) CreateContainer(hostname string, image string) (bool, error) {
	// Container creation request
	req := api.ContainersPost{
		Name: hostname,
		Source: api.ContainerSource{
			Type:     "image",
			Server:   "https://cloud-images.ubuntu.com/releases",
			Protocol: "simplestreams",
			Alias:    image,
		},
	}

	// Get LXD to create the container (background operation)
	op, err := l.conn.CreateContainer(req)
	if err != nil {
		return false, err
	}

	// Wait for the operation to complete
	err = op.Wait()
	if err != nil {
		return false, err
	}

	// Get LXD to start the container (background operation)
	reqState := api.ContainerStatePut{
		Action:  "start",
		Timeout: -1,
	}

	op, err = l.conn.UpdateContainerState(hostname, reqState, "")
	if err != nil {
		return false, err
	}

	// Wait for the operation to complete
	err = op.Wait()
	if err != nil {
		return false, err
	}

	return true, nil
}
