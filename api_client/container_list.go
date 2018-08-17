package api_client

import (
	"encoding/json"
)

type ContainerList struct {
	ApiVersion int               `json:"api_version"`
	Data       ContainerListData `json:"data"`
}

type ContainerListData struct {
	Containers []Container `json:"containers"`
}

func NewContainerListFromByte(b []byte) (*ContainerList, error) {
	var c ContainerList
	err := json.Unmarshal(b, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
