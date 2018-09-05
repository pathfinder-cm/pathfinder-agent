package daemon

import (
	"errors"
	"time"

	client "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
	"github.com/pathfinder-cm/pathfinder-go-client/pfmodel"
)

type apiContainers []api.Container

type LXD struct {
	hostname  string
	localSrv  client.ContainerServer
	targetSrv client.ContainerServer
}

func (a apiContainers) toContainerList() *pfmodel.ContainerList {
	containerList := make(pfmodel.ContainerList, len(a))
	for i, c := range a {
		containerList[i] = pfmodel.Container{
			Hostname: c.Name,
		}
	}
	return &containerList
}

func NewLXD(hostname string, socketPath string) (*LXD, error) {
	localSrv, err := client.ConnectLXDUnix(socketPath, nil)
	if err != nil {
		return nil, err
	}
	targetSrv := localSrv.UseTarget(hostname)

	return &LXD{
		hostname:  hostname,
		localSrv:  localSrv,
		targetSrv: targetSrv,
	}, nil
}

func (l *LXD) ListContainers() (*pfmodel.ContainerList, error) {
	res, err := l.targetSrv.GetContainers()
	if err != nil {
		return nil, err
	}

	containerList := apiContainers(res).toContainerList()

	return containerList, nil
}

func (l *LXD) CreateContainer(hostname string, image string) (bool, string, error) {
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
	// but check first if cluster mode is active
	var err error
	var op client.Operation
	if l.targetSrv.IsClustered() {
		op, err = l.targetSrv.CreateContainer(req)
	} else {
		op, err = l.localSrv.CreateContainer(req)
	}

	if err != nil {
		return false, "", err
	}

	// Wait for the operation to complete
	err = op.Wait()
	if err != nil {
		return false, "", err
	}

	// Get LXD to start the container (background operation)
	startReq := api.ContainerStatePut{
		Action:  "start",
		Timeout: -1,
	}

	op, err = l.targetSrv.UpdateContainerState(hostname, startReq, "")
	if err != nil {
		return false, "", err
	}

	// Wait for the operation to complete
	err = op.Wait()
	if err != nil {
		return false, "", err
	}

	// Wait for ipaddress to be available
	ipaddress := ""
	found := false
	timeLimit := time.Now().Add(60 * time.Second)

	for !found && time.Now().Before(timeLimit) {
		state, _, err := l.targetSrv.GetContainerState(hostname)
		if err != nil {
			return false, "", err
		}

		addresses := state.Network["eth0"].Addresses
		for _, address := range addresses {
			if address.Family == "inet" {
				ipaddress = address.Address
			}
		}

		if ipaddress != "" {
			found = true
			break
		}

		time.Sleep(time.Duration(3) * time.Second)
	}

	if !found {
		return false, "", errors.New("Missing ip address")
	}

	return true, ipaddress, nil
}

func (l *LXD) DeleteContainer(hostname string) (bool, error) {
	// Get LXD to stop the container (background operation)
	stopReq := api.ContainerStatePut{
		Action:  "stop",
		Timeout: 60,
	}

	// Stop the container
	// We don't care if container already stopped
	op, _ := l.targetSrv.UpdateContainerState(hostname, stopReq, "")
	op.Wait()

	// Get LXD to delete the container (background operation)
	op, err := l.targetSrv.DeleteContainer(hostname)
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
