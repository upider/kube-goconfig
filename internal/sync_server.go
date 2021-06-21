package internal

import (
	"context"
	"kube-goconfig/pkg"
	"sync"

	mapset "github.com/deckarep/golang-set"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Server interface {
	Start(signalController *pkg.SignalController)
	Stop()
}

// SyncServer 同步nacos配置到k8s
type SyncServer struct {
	NacosClients map[config_client.IConfigClient]mapset.Set `json:"nacosClients"`
	K8sClientset *kubernetes.Clientset
	SyncConfig   *SyncConfiguration
	mutex        sync.Mutex
}

func (s *SyncServer) Start(signalController *pkg.SignalController) {
	log.Info("SyncServer start...\n")
	ctx, cancel := context.WithCancel(context.Background())
	cfg := s.SyncConfig

	// sync k8s namespace
	if cfg.AutoCreatek8sNs {
		for _, namesapce := range cfg.SyncNamespaces {
			ns := &corev1.Namespace{}
			_, err := s.K8sClientset.CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
			if err != nil {
				log.Error(err)
			} else {
				log.Infof("%s namespace already exists", namesapce)
			}
		}
	}

	for {
		select {
		case <-signalController.SignalChan:
			s.Stop()
			cancel()
		default:
			// sync config
			for nacosClient, configs := range s.NacosClients {
				newConfigSet := mapset.NewSet()

				for i := 1; ; i++ {
					page, err := nacosClient.SearchConfig(vo.SearchConfigParam{
						Search:   "blur",
						DataId:   "",
						Group:    "",
						PageNo:   i,
						PageSize: 10,
					})
					if err != nil {
						continue
					}
					for _, config := range page.PageItems {
						newConfigSet.Add(config)
					}

					if page.PagesAvailable == 0 {
						break
					}
				}

				if configs.Equal(newConfigSet) {
					// 如果没有变化则不添加新的ListenConfig，也不删除ListenConfig
					continue
				}

				// 添加新的ListenConfig或删除ListenConfig
				interSet := configs.Intersect(newConfigSet)
				deleteListenSet := configs.Difference(interSet)
				addListenSet := configs.Difference(newConfigSet)

				delIt := deleteListenSet.Iterator()
				defer delIt.Stop()

				for configItem := range delIt.C {
					nacosClient.CancelListenConfig(vo.ConfigParam{
						DataId: configItem.(model.ConfigItem).DataId,
						Group:  configItem.(model.ConfigItem).Group,
					})
				}

				addIt := addListenSet.Iterator()
				immutable := false
				for configItem := range addIt.C {
					nacosClient.ListenConfig(vo.ConfigParam{
						DataId: configItem.(model.ConfigItem).DataId,
						Group:  configItem.(model.ConfigItem).Group,
						OnChange: func(namespace, group, dataId, data string) {
							configMap := &corev1.ConfigMap{
								ObjectMeta: metav1.ObjectMeta{
									Name:      dataId,
									Namespace: namespace,
									Labels:    map[string]string{"group": group},
								},
								Immutable: &immutable,
								Data:      map[string]string{dataId: data},
							}
							s.mutex.Lock()
							s.K8sClientset.CoreV1().ConfigMaps(namespace).Create(ctx, configMap, metav1.CreateOptions{})
							s.mutex.Unlock()
						},
					})
				}

				s.NacosClients[nacosClient] = newConfigSet
			}
		}
	}
}

func (s *SyncServer) Stop() {
	log.Info("SyncServer stop...\n")
}

func NewSyncServer(cfg SyncConfiguration) (Server, error) {
	var syncServer SyncServer

	// creates the in-cluster config
	config, err := pkg.GetKubeConfig()
	if err != nil {
		return nil, err
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	syncServer.K8sClientset = clientset

	// create nacos client
	syncServer.NacosClients = make(map[config_client.IConfigClient]mapset.Set, len(cfg.SyncNamespaces))

	serverConfigs := make([]constant.ServerConfig, len(cfg.NacosIPs))
	for _, ip := range cfg.NacosIPs {
		serverConfigs = append(serverConfigs, *constant.NewServerConfig(
			ip,
			cfg.NacosPort,
			constant.WithScheme("http"),
			constant.WithContextPath("/nacos"),
		))
	}

	for _, namespace := range cfg.SyncNamespaces {
		clientConfig := *constant.NewClientConfig(
			constant.WithNamespaceId(namespace), //When namespace is public, fill in the blank string here.
			constant.WithTimeoutMs(5000),
			constant.WithNotLoadCacheAtStart(true),
			constant.WithLogDir("/tmp/nacos/log"),
			constant.WithCacheDir("/tmp/nacos/cache"),
			constant.WithRotateTime("1h"),
			constant.WithMaxAge(3),
			constant.WithLogLevel(cfg.LogLevel),
		)

		configClient, err := clients.NewConfigClient(
			vo.NacosClientParam{
				ClientConfig:  &clientConfig,
				ServerConfigs: serverConfigs,
			},
		)

		if err != nil {
			return nil, err
		}

		syncServer.NacosClients[configClient] = mapset.NewSet()
	}

	return &syncServer, nil
}
