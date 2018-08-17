package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/lxc/lxd/shared/api"
)

func findContainer(containers []api.Container, name string) int {
	for i, c := range containers {
		if c.Name == name {
			return i
		}
	}

	return -1
}

func main() {
	for {
		// Get from API Server
		b, _ := ioutil.ReadFile("/opt/projects/golang/src/github.com/giosakti/pathfinder-agent/fixtures/scheduled-containers.json")

		scheduled, _ := NewContainersFromByte(b)
		containers := scheduled.Data.Containers

		// Get from LXC Host
		local, _ := listContainers()

		// Compare API Server and LXC Host
		for _, c := range containers {
			j := findContainer(local, c.Name)
			if j == -1 {
				fmt.Println("Creating Container", c.Name)
				createContainer(c.Name)
			} else {
				local = append(local[:j], local[j+1:]...)
			}
		}

		time.Sleep(5 * time.Second)
	}
}
