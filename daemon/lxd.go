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

func (a apiContainers) toContainerList() *model.ContainerList {
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

	containerList := apiContainers(res).toContainerList()

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
	startReq := api.ContainerStatePut{
		Action:  "start",
		Timeout: -1,
	}

	op, err = l.conn.UpdateContainerState(hostname, startReq, "")
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

func (l *LXD) DeleteContainer(hostname string) (bool, error) {
	// Get LXD to stop the container (background operation)
	stopReq := api.ContainerStatePut{
		Action:  "stop",
		Timeout: 60,
	}

	// Stop the container
	// We don't care if container already stopped
	op, _ := l.conn.UpdateContainerState(hostname, stopReq, "")
	op.Wait()

	// Get LXD to delete the container (background operation)
	op, err := l.conn.DeleteContainer(hostname)
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
