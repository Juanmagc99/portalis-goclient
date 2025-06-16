package pclient

type Instance struct {
	ServiceName string            `json:"serviceName"`
	InstanceID  string            `json:"instanceID"`
	Host        string            `json:"host"`
	Port        int               `json:"port"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}
