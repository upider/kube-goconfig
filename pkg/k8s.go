package pkg

import (
	"path/filepath"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

//GetKubeConfig get kubernetes config
func GetKubeConfig() (*rest.Config, error) {
	var kubeConfig *rest.Config
	var err error

	// creates the in-cluster config
	kubeConfig, err = rest.InClusterConfig()
	if err == nil {
		return kubeConfig, nil
	} else if err == rest.ErrNotInCluster {
		// creates the out-cluster config
		if home := homedir.HomeDir(); home != "" {
			config := filepath.Join(home, ".kube", "config")
			kubeConfig, err = clientcmd.BuildConfigFromFlags("", config)
			if err != nil {
				panic(err.Error())
			}
		} else {
			panic(err.Error())
		}
	} else {
		return nil, err
	}

	return kubeConfig, nil
}
