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
	MarkContainerAsProvisionError(node string, hostname string) (bool, error)
	MarkContainerAsDeleted(node string, hostname string) (bool, error)
}

type pfclient struct {
	cluster      string
	httpClient   *http.Client
	pfServerAddr string
	pfApiPath    map[string]string
}

func NewPfclient(
	cluster string,
	httpClient *http.Client,
	pfServerAddr string,
	pfApiPath map[string]string) Pfclient {

	return &pfclient{
		cluster:      cluster,
		httpClient:   httpClient,
		pfServerAddr: pfServerAddr,
		pfApiPath:    pfApiPath,
	}
}

func (p *pfclient) FetchContainersFromServer(node string) (*model.ContainerList, error) {
	address := fmt.Sprintf("%s/%s", p.pfServerAddr, p.pfApiPath["ListContainers"])
	q := url.Values{}
	q.Add("cluster_name", p.cluster)
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
	address := fmt.Sprintf("%s/%s", p.pfServerAddr, p.pfApiPath["MarkProvisioned"])
	form := url.Values{}
	form.Set("cluster_name", p.cluster)
	form.Set("node_hostname", node)
	form.Set("hostname", hostname)
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

func (p *pfclient) MarkContainerAsProvisionError(node string, hostname string) (bool, error) {
	address := fmt.Sprintf("%s/%s", p.pfServerAddr, p.pfApiPath["MarkProvisionError"])
	form := url.Values{}
	form.Set("cluster_name", p.cluster)
	form.Set("node_hostname", node)
	form.Set("hostname", hostname)
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

func (p *pfclient) MarkContainerAsDeleted(node string, hostname string) (bool, error) {
	address := fmt.Sprintf("%s/%s", p.pfServerAddr, p.pfApiPath["MarkDeleted"])
	form := url.Values{}
	form.Set("cluster_name", p.cluster)
	form.Set("node_hostname", node)
	form.Set("hostname", hostname)
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
