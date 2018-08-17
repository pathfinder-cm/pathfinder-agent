package main

import (
	"fmt"
	"io/ioutil"
	"time"

	. "github.com/giosakti/pathfinder-agent/api_client"
	. "github.com/giosakti/pathfinder-agent/lxd_client"
)

func main() {
	for {
		// Get from API Server
		b, _ := ioutil.ReadFile("/opt/projects/golang/src/github.com/giosakti/pathfinder-agent/fixtures/scheduled-containers.json")

		scheduled, _ := NewContainerListFromByte(b)
		containers := scheduled.Data.Containers

		// Get from LXC Host
		local, _ := ListContainers()

		// Compare API Server and LXC Host
		for _, c := range containers {
			j := FindContainer(local, c.Name)
			if j == -1 {
				fmt.Println("Creating Container", c.Name)
				CreateContainer(c.Name)
			} else {
				local = append(local[:j], local[j+1:]...)
			}
		}

		time.Sleep(5 * time.Second)
	}
}
