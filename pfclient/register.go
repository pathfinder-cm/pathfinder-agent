package pfclient

import (
	"encoding/json"
)

type RegisterRes struct {
	ApiVersion string          `json:"api_version"`
	Data       RegisterDataRes `json:"data"`
}

type RegisterDataRes struct {
	Hostname            string `json:"hostname"`
	AuthenticationToken string `json:"authentication_token"`
}

func NewRegisterFromByte(b []byte) (*RegisterDataRes, error) {
	var res RegisterRes
	err := json.Unmarshal(b, &res)
	if err != nil {
		return nil, err
	}

	return &res.Data, nil
}
