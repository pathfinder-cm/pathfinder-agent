package lxd_client

import (
	lxd "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
)

func ListContainers() ([]api.Container, error) {
	// Connect to LXD over the Unix socket
	c, err := lxd.ConnectLXDUnix("/var/snap/lxd/common/lxd/unix.socket", nil)
	if err != nil {
		return nil, err
	}

	containers, err := c.GetContainers()
	if err != nil {
		return nil, err
	}

	return containers, nil
}

func FindContainer(containers []api.Container, name string) int {
	for i, c := range containers {
		if c.Name == name {
			return i
		}
	}

	return -1
}

func CreateContainer(name string) (bool, error) {
	// Connect to LXD over the Unix socket
	c, err := lxd.ConnectLXDUnix("/var/snap/lxd/common/lxd/unix.socket", nil)
	if err != nil {
		return false, err
	}

	// Container creation request
	req := api.ContainersPost{
		Name: name,
		Source: api.ContainerSource{
			Type:     "image",
			Server:   "https://cloud-images.ubuntu.com/releases",
			Protocol: "simplestreams",
			Alias:    "16.04",
		},
	}

	// Get LXD to create the container (background operation)
	op, err := c.CreateContainer(req)
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

	op, err = c.UpdateContainerState(name, reqState, "")
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
