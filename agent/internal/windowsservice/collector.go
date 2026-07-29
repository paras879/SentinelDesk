package windowsservice


type ServiceInfo struct {
	ServiceName    string `json:"service_name"`
	DisplayName    string `json:"display_name"`
	Status         string `json:"status"`
	StartType      string `json:"start_type"`
	ExecutablePath string `json:"executable_path"`
	PID            int32  `json:"pid"`
	ServiceAccount string `json:"service_account"`
	Description    string `json:"description"`
	CanStop        bool   `json:"can_stop"`
	CanPause       bool   `json:"can_pause"`
	AcceptShutdown bool   `json:"accept_shutdown"`
}

