package main

import (
	"encoding/json"
)

type ContainersResponse struct {
	ApiVersion int                    `json:"api_version"`
	Data       ContainersResponseData `json:"data"`
}

type ContainersResponseData struct {
	Containers []Container `json:"containers"`
}

type Container struct {
	Name string `json:"name"`
}

func NewContainersFromByte(b []byte) (*ContainersResponse, error) {
	var c ContainersResponse
	err := json.Unmarshal(b, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
