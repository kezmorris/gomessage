package main

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// KubeClient interracts with the kubenetes cluster
type KubeClient struct {
	clientSet *kubernetes.Clientset
}

// NewKubeClient creates and connectes a KubeClient object
func NewKubeClient() (*KubeClient, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("Error getting cluster config: %v", err)
	}
	newclientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("Error creating clientset: %v", err)
	}
	return &KubeClient{newclientset}, nil
}
