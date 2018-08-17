package agent

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/giosakti/pathfinder-agent/lxd_client"
	"github.com/giosakti/pathfinder-agent/pfclient"
)

type Agent interface {
	Run()
}

type agent struct {
}

func NewAgent() Agent {
	return &agent{}
}

func (a *agent) Run() {
	for {
		// Get from API Server
		b, _ := ioutil.ReadFile("/opt/projects/golang/src/github.com/giosakti/pathfinder-agent/fixtures/scheduled-containers.json")
		rc, _ := pfclient.NewContainerListFromByte(b)

		// Get from local daemon
		lc, _ := lxd_client.ListContainers()

		// Compare containers from server and local daemon
		for _, c := range *rc {
			i := lxd_client.FindContainer(lc, c.Name)
			if i == -1 {
				fmt.Println("Creating Container", c.Name)
				lxd_client.CreateContainer(c.Name)
			} else {
				lc = append(lc[:i], lc[i+1:]...)
			}
		}

		time.Sleep(5 * time.Second)
	}
}
