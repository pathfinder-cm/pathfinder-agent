package pfclient

import (
	"encoding/json"

	"github.com/giosakti/pathfinder-agent/model"
)

type ContainerListRes struct {
	ApiVersion int                  `json:"api_version"`
	Data       ContainerListDataRes `json:"data"`
}

type ContainerListDataRes struct {
	Containers []ContainerRes `json:"containers"`
}

type ContainerRes struct {
	Name string `json:"name"`
}

func NewContainerListFromByte(b []byte) (*model.ContainerList, error) {
	res, err := newContainerListResFromByte(b)
	if err != nil {
		return nil, err
	}

	cl := make(model.ContainerList, len(res.Data.Containers))
	for i, c := range res.Data.Containers {
		cl[i] = model.Container{
			Name: c.Name,
		}
	}

	return &cl, nil
}

func newContainerListResFromByte(b []byte) (*ContainerListRes, error) {
	var res ContainerListRes
	err := json.Unmarshal(b, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
