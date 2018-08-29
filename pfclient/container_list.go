package pfclient

import (
	"encoding/json"

	"github.com/pathfinder-cm/pathfinder-agent/model"
)

type ContainerListRes struct {
	ApiVersion string               `json:"api_version"`
	Data       ContainerListDataRes `json:"data"`
}

type ContainerListDataRes struct {
	Items []ContainerRes `json:"items"`
}

type ContainerRes struct {
	Hostname string `json:"hostname"`
	Image    string `json:"image"`
	Status   string `json:"status"`
}

func NewContainerListFromByte(b []byte) (*model.ContainerList, error) {
	var res ContainerListRes
	err := json.Unmarshal(b, &res)
	if err != nil {
		return nil, err
	}

	cl := make(model.ContainerList, len(res.Data.Items))
	for i, c := range res.Data.Items {
		cl[i] = model.Container{
			Hostname: c.Hostname,
			Image:    c.Image,
			Status:   c.Status,
		}
	}

	return &cl, nil
}
