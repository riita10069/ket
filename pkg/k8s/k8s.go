package k8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type ClientGo struct {
	KubeconfigPath string
	ClientSet      *kubernetes.Clientset
}

func NewClientGo(kubeConfigPath string) (*ClientGo, error) {
	context, err := currentContext(kubeConfigPath)
	if err != nil {
		return nil, err
	}

	clientSet, err := createClient(context)
	if err != nil {
		return nil, err
	}

	return &ClientGo{
		KubeconfigPath: kubeConfigPath,
		ClientSet:      clientSet,
	}, nil
}

// currentContext use the current context in kubeconfig
func currentContext(kubeConfig string) (*rest.Config, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// createClient create the clientset
func createClient(config *rest.Config) (*kubernetes.Clientset, error) {
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientSet, nil
}
