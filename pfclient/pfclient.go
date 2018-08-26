package pfclient

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/giosakti/pathfinder-agent/model"
	log "github.com/sirupsen/logrus"
)

type Pfclient interface {
	FetchContainersFromServer(node string) (*model.ContainerList, error)
	MarkContainerAsProvisioned(node string, hostname string) (bool, error)
}

type pfclient struct {
	httpClient           *http.Client
	pfServerAddr         string
	pfListContainersPath string
	pfProvisionedPath    string
}

func NewPfclient(
	httpClient *http.Client,
	pfServerAddr string,
	pfListContainersPath string,
	pfProvisionedPath string) Pfclient {

	return &pfclient{
		httpClient:           httpClient,
		pfServerAddr:         pfServerAddr,
		pfListContainersPath: pfListContainersPath,
		pfProvisionedPath:    pfProvisionedPath,
	}
}

func (p *pfclient) FetchContainersFromServer(node string) (*model.ContainerList, error) {
	address := fmt.Sprintf("%s/%s", p.pfServerAddr, p.pfListContainersPath)
	q := url.Values{}
	q.Add("node_hostname", node)

	req, _ := http.NewRequest(http.MethodGet, address, nil)
	req.URL.RawQuery = q.Encode()

	res, err := p.httpClient.Do(req)
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
	serverContainers, err := NewContainerListFromByte(b)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return serverContainers, nil
}

func (p *pfclient) MarkContainerAsProvisioned(node string, hostname string) (bool, error) {
	address := fmt.Sprintf("%s/%s", p.pfServerAddr, p.pfProvisionedPath)
	form := url.Values{}
	form.Set("node_hostname", node)
	form.Add("hostname", hostname)
	body := bytes.NewBufferString(form.Encode())

	req, err := http.NewRequest(http.MethodPost, address, body)

	res, err := p.httpClient.Do(req)
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
