package pclient

import "time"

type Instance struct {
	ServiceName string            `json:"serviceName" validate:"required"`
	InstanceID  string            `json:"instanceID"  validate:"required"`
	Host        string            `json:"host"        validate:"required,hostname|ip"`
	Port        int               `json:"port"        validate:"required,min=1"`
	Metadata    map[string]string `json:"metadata"    validate:"required"`
	LastSeen    time.Time         `json:"lastSeen"    validate:"required"`
}
