package pfclient

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pathfinder-cm/pathfinder-agent/model"
	log "github.com/sirupsen/logrus"
)

type Pfclient interface {
	FetchContainersFromServer(node string) (*model.ContainerList, error)
	UpdateIpaddress(node string, hostname string, ipaddress string) (bool, error)
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
	addr := fmt.Sprintf("%s/%s", p.pfServerAddr, p.pfApiPath["ListContainers"])
	u, err := url.Parse(addr)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	q := u.Query()
	q.Set("cluster_name", p.cluster)
	q.Set("node_hostname", node)
	u.RawQuery = q.Encode()

	req, _ := http.NewRequest(http.MethodGet, u.String(), nil)
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

func (p *pfclient) UpdateIpaddress(node string, hostname string, ipaddress string) (bool, error) {
	addr := fmt.Sprintf("%s/%s", p.pfServerAddr, p.pfApiPath["UpdateIpaddress"])
	u, err := url.Parse(addr)
	if err != nil {
		log.Error(err.Error())
		return false, err
	}
	q := u.Query()
	q.Set("cluster_name", p.cluster)
	q.Set("node_hostname", node)
	q.Set("hostname", hostname)
	u.RawQuery = q.Encode()

	form := url.Values{}
	form.Set("ipaddress", ipaddress)
	body := bytes.NewBufferString(form.Encode())

	req, err := http.NewRequest(http.MethodPost, u.String(), body)
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

func (p *pfclient) MarkContainerAsProvisioned(node string, hostname string) (bool, error) {
	addr := fmt.Sprintf("%s/%s", p.pfServerAddr, p.pfApiPath["MarkProvisioned"])
	u, err := url.Parse(addr)
	if err != nil {
		log.Error(err.Error())
		return false, err
	}
	q := u.Query()
	q.Set("cluster_name", p.cluster)
	q.Set("node_hostname", node)
	q.Set("hostname", hostname)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodPost, u.String(), nil)
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
	addr := fmt.Sprintf("%s/%s", p.pfServerAddr, p.pfApiPath["MarkProvisionError"])
	u, err := url.Parse(addr)
	if err != nil {
		log.Error(err.Error())
		return false, err
	}
	q := u.Query()
	q.Set("cluster_name", p.cluster)
	q.Set("node_hostname", node)
	q.Set("hostname", hostname)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodPost, u.String(), nil)

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
	addr := fmt.Sprintf("%s/%s", p.pfServerAddr, p.pfApiPath["MarkDeleted"])
	u, err := url.Parse(addr)
	if err != nil {
		log.Error(err.Error())
		return false, err
	}
	q := u.Query()
	q.Set("cluster_name", p.cluster)
	q.Set("node_hostname", node)
	q.Set("hostname", hostname)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodPost, u.String(), nil)

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
