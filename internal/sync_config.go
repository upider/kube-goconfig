package internal

type SyncConfiguration struct {
	SyncNamespaces  []string `json:"syncNamespaces"`
	AutoCreatek8sNs bool     `json:"autoCreatek8sNs"`
	LogLevel        string   `json:"logLevel"`
	NacosIPs        []string `json:"nacosIPs"`
	NacosPort       uint64   `json:"nacosPort"`
	ConfigScanTime  int64    `json:"configScanTime"`
}
