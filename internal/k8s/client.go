package k8s

import (
	"fmt"

	"github.com/maz/k8s-mini-explorer/internal/config"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func NewClient(cfg *config.Config) (kubernetes.Interface, error) {
	var restCfg *rest.Config
	var err error

	if cfg.InCluster {
		restCfg, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("in-cluster config: %w", err)
		}
	} else {
		loadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: cfg.KubeconfigPath}
		overrides := &clientcmd.ConfigOverrides{}
		if cfg.KubeContext != "" {
			overrides.CurrentContext = cfg.KubeContext
		}
		clientCfg := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, overrides)
		restCfg, err = clientCfg.ClientConfig()
		if err != nil {
			return nil, fmt.Errorf("kubeconfig: %w", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(restCfg)
	if err != nil {
		return nil, fmt.Errorf("kubernetes client: %w", err)
	}

	return clientset, nil
}
