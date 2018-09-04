package model

type Metrics struct {
	Memory Memory `json:"memory"`
}

type Memory struct {
	Used  uint64 `json:"used"`
	Free  uint64 `json:"free"`
	Total uint64 `json:"total"`
}
