package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port                  int
	KubeconfigPath        string
	KubeContext           string
	DefaultNS             string
	InCluster             bool
	InsecureSkipTLSVerify bool
}

func Load() *Config {
	port := 8080
	if p := os.Getenv("PORT"); p != "" {
		if v, err := strconv.Atoi(p); err == nil {
			port = v
		}
	}

	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		if home, err := os.UserHomeDir(); err == nil {
			kubeconfig = home + "/.kube/config"
		}
	}

	inCluster := os.Getenv("IN_CLUSTER") == "true"

	defaultNS := os.Getenv("K8S_NAMESPACE")
	if defaultNS == "" {
		defaultNS = "default"
	}

	return &Config{
		Port:                  port,
		KubeconfigPath:        kubeconfig,
		KubeContext:           os.Getenv("K8S_CONTEXT"),
		DefaultNS:             defaultNS,
		InCluster:             inCluster,
		InsecureSkipTLSVerify: os.Getenv("INSECURE_SKIP_TLS_VERIFY") == "true",
	}
}
