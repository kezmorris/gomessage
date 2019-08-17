package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// KubeClient interracts with the kubenetes cluster
type KubeClient struct {
	clientset *kubernetes.Clientset
	logger    *Logger
}

// NewInClusterKubeClient creates and connectes a KubeClient object.
// Configured inside a k8s cluster
func NewInClusterKubeClient(logger *Logger) (*KubeClient, error) {
	config, err := rest.InClusterConfig() // NEED A LOCAL CONTEXT CONFIG
	if err != nil {
		return nil, fmt.Errorf("Error getting cluster config: %v", err)
	}
	newclientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("Error creating clientset: %v", err)
	}
	return &KubeClient{newclientset, logger}, nil
}

// NewKubeClient creates and connectes a KubeClient object.
// Configured outside a k8s cluster using local kubecfg.
func NewKubeClient(logger *Logger) (*KubeClient, error) {
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
	return &KubeClient{newclientset, logger}, nil
}

// Run begins the kubernetes interfacing
func (kc *KubeClient) Run() {
	go kc.monitorPods()
}

func (kc *KubeClient) monitorPods() {
	for {
		var namespace string

		if os.Getenv("NAMESPACE") != "" {
			namespace = os.Getenv("NAMESPACE")
		} else {
			namespace = "" // Either the deployment is misconfigured or app is not running in cluster
		}
		for {
			pods, err := kc.clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{})
			if err != nil {
				kc.logger.Log(fmt.Sprintf("Error calling clientset: %v", err))
			}
			kc.logger.Log(fmt.Sprintf("There are %d pods\n", len(pods.Items)))
			time.Sleep(20 * time.Second)
		}
	}
}

// RequestPort interfaces with the kubeclient main process to get a port to allocate to the pod
func (kc *KubeClient) RequestPort() int {
	var port int
	port = 30001
	kc.logger.Log(fmt.Sprintf("Got port: %d", port))
	return port
}
