package models

type Step struct {
	Name         string `json:"name"`
	EndpointType string `json:"endpoint_type"`
	EndpointID   int64  `json:"endpoint_id"`
}
