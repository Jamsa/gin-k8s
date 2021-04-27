package k8s

import (
	"flag"
	"path/filepath"
	"sync"

	"github.com/jamsa/gin-k8s/pkg/logging"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var k8sClients = &sync.Map{} //请求map

func GetClient(clusterID string) (*kubernetes.Clientset, error) {
	client, ok := k8sClients.Load(clusterID)
	if ok {
		logging.Info("从缓存中得到", clusterID, "的连接")
		return client.(*kubernetes.Clientset), nil
	}

	logging.Info("为", clusterID, "创建连接")
	//TODO: 从数据库读取配置
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}
	//clientcmd.NewClientConfigFromBytes()

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	logging.Info(clusterID, "连接创建成功")
	k8sClients.Store(clusterID, clientset)
	return clientset, nil
}
