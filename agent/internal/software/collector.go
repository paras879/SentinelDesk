package software


type Software struct {
	Name            string `json:"name"`
	Version         string `json:"version"`
	Publisher       string `json:"publisher"`
	InstallDate     string `json:"install_date"`
	InstallLocation string `json:"install_location"`
	EstimatedSize   int64  `json:"estimated_size"`
}


