package agent

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/giosakti/pathfinder-agent/daemon"
	"github.com/giosakti/pathfinder-agent/model"
	"github.com/giosakti/pathfinder-agent/pfclient"
	log "github.com/sirupsen/logrus"
)

type Agent interface {
	Run()
}

type agent struct {
	nodeHostname         string
	containerDaemon      daemon.ContainerDaemon
	httpClient           *http.Client
	pfServerAddr         string
	pfListContainersPath string
	pfProvisionedPath    string
}

func NewAgent(
	nodeHostname string,
	containerDaemon daemon.ContainerDaemon,
	httpClient *http.Client,
	pfServerAddr string,
	pfListContainersPath string,
	pfProvisionedPath string) Agent {

	return &agent{
		nodeHostname:         nodeHostname,
		containerDaemon:      containerDaemon,
		httpClient:           httpClient,
		pfServerAddr:         pfServerAddr,
		pfListContainersPath: pfListContainersPath,
		pfProvisionedPath:    pfProvisionedPath,
	}
}

func (a *agent) Run() {
	for {
		// Add delay between processing
		time.Sleep(5 * time.Second)

		serverContainers, err := fetchContainersFromServer(
			a.httpClient,
			a.pfServerAddr,
			a.pfListContainersPath,
			a.nodeHostname,
		)
		if err != nil {
			continue
		}

		localContainers, err := a.containerDaemon.ListContainers()
		if err != nil {
			continue
		}

		// Compare containers between server and local daemon
		// Do action as necessary
		for _, sc := range *serverContainers {
			ok, _ := a.provisionContainer(sc, localContainers)
			if !ok {
				continue
			}
		}
	}
}

func (a *agent) provisionContainer(sc model.Container, localContainers *model.ContainerList) (bool, error) {
	i := localContainers.FindByHostname(sc.Hostname)
	if i == -1 {
		log.WithFields(log.Fields{
			"hostname": sc.Hostname,
			"number":   sc.Image,
		}).Info("Creating container")
		a.containerDaemon.CreateContainer(sc.Hostname, sc.Image)

		ok, err := markContainerAsProvisioned(
			a.httpClient,
			a.pfServerAddr,
			a.pfProvisionedPath,
			a.nodeHostname,
			sc.Hostname,
		)
		if !ok {
			return false, err
		}

		log.WithFields(log.Fields{
			"hostname": sc.Hostname,
			"number":   sc.Image,
		}).Info("Container created")
	} else {
		localContainers.DeleteAt(i)
		log.WithFields(log.Fields{
			"hostname": sc.Hostname,
			"number":   sc.Image,
		}).Info("Container already exist")
	}

	return true, nil
}

func fetchContainersFromServer(client *http.Client, addr string, path string, node string) (*model.ContainerList, error) {
	address := fmt.Sprintf("%s/%s", addr, path)
	q := url.Values{}
	q.Add("node_hostname", node)

	req, _ := http.NewRequest("GET", address, nil)
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(res.Body)
		log.Error(string(b))
		return nil, errors.New(string(b))
	}

	b, _ := ioutil.ReadAll(res.Body)
	serverContainers, err := pfclient.NewContainerListFromByte(b)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return serverContainers, nil
}

func markContainerAsProvisioned(client *http.Client, addr string, path string, node string, hostname string) (bool, error) {
	address := fmt.Sprintf("%s/%s", addr, path)
	form := url.Values{}
	form.Set("node_hostname", node)
	form.Add("hostname", hostname)
	body := bytes.NewBufferString(form.Encode())

	req, err := http.NewRequest("POST", address, body)

	res, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return false, err
	}

	if res.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(res.Body)
		log.Error(string(b))
		return false, errors.New(string(b))
	}

	return true, nil
}
