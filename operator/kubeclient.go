package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// KubeClient interracts with the kubenetes cluster
type KubeClient struct {
	clientset *kubernetes.Clientset
}

// NewInClusterKubeClient creates and connectes a KubeClient object.
// Configured inside a k8s cluster
func NewInClusterKubeClient() (*KubeClient, error) {
	config, err := rest.InClusterConfig() // NEED A LOCAL CONTEXT CONFIG
	if err != nil {
		return nil, fmt.Errorf("Error getting cluster config: %v", err)
	}
	newclientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("Error creating clientset: %v", err)
	}
	return &KubeClient{newclientset}, nil
}

// NewKubeClient creates and connectes a KubeClient object.
// Configured outside a k8s cluster using local kubecfg.
func NewKubeClient() (*KubeClient, error) {
	var kubeconfig *string

	homeDir := os.Getenv("HOME")
	kubeconfig = flag.String("kubeconfig", filepath.Join(homeDir, ".kube", "config"), "absolute path to the kubeconfig file")
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("Error building client: %v", err)
	}
	newclientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("Error creating clientset: %v", err)
	}
	return &KubeClient{newclientset}, nil
}
