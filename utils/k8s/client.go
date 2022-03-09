package k8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func KubernetesClient() (*kubernetes.Clientset, *rest.Config, string, error) {
	var err error
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)
	namespace, _, err := kubeConfig.Namespace()
	if err != nil {
		return nil, nil, "", err
	}
	restconfig, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, nil, "", err
	}
	clientset, err := kubernetes.NewForConfig(restconfig)
	if err != nil {
		return nil, nil, "", err
	}
	return clientset, restconfig, namespace, nil
}
