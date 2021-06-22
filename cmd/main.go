package main

import (
	"time"

	"kube-goconfig/internal"
	"kube-goconfig/pkg"

	flag "github.com/spf13/pflag"
)

var (
	namespaces      []string
	nacosIPs        []string
	nacosPort       uint64
	logLevel        string
	autoCreatek8sNs bool
	configScanTime  int64
)

func main() {
	flag.StringSliceVar(&namespaces, "namespaces", nil, "namespaces need to synchronization configuration")
	flag.StringSliceVar(&nacosIPs, "nacosIPs", nil, "nacos servers ips")
	flag.Uint64Var(&nacosPort, "nacosPort", 8848, "nacos server port")
	flag.Int64Var(&configScanTime, "configScanTime", 10, "synchronization interval")
	flag.StringVar(&logLevel, "logLevel", "info", "log level")
	flag.BoolVar(&autoCreatek8sNs, "autoCreatek8sNs", false, "auto create kubernetes namespace (default false)")
	flag.Parse()

	if nacosIPs == nil || namespaces == nil {
		flag.Usage()
		return
	}

	syncConfig := internal.SyncConfiguration{
		SyncNamespaces:  namespaces,
		NacosIPs:        nacosIPs,
		NacosPort:       nacosPort,
		LogLevel:        logLevel,
		AutoCreatek8sNs: autoCreatek8sNs,
		ConfigScanTime:  configScanTime,
	}

	server, err := internal.NewSyncServer(&syncConfig)
	if err != nil {
		panic(err)
	}

	signalController := pkg.NewSignalController(time.Second * 5)

	server.Start(signalController)
}
