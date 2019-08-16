package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// KubeClient interracts with the kubenetes cluster
type KubeClient struct {
	clientset *kubernetes.Clientset
}

// NewKubeClient creates and connectes a KubeClient object.
func NewKubeClient() (*KubeClient, error) {
	var kubeconfig *string

	homeDir := os.Getenv("HOME")
	kubeconfig = flag.String("kubeconfig", filepath.Join(homeDir, ".kube", "config"), "absolute path to the kubeconfig file")
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	newclientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("Error creating clientset: %v", err)
	}
	return &KubeClient{newclientset}, nil
}

func main() {
	client, err := NewKubeClient()
	if err != nil {
		fmt.Printf("Error getting kubeclient: %v", err)
		panic(err)
	}
	for {
		pods, err := client.clientset.CoreV1().Pods("").List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

		// Examples for error handling:
		// - Use helper functions like e.g. errors.IsNotFound()
		// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
		namespace := "default"
		pod := "example-xxxxx"
		_, err = client.clientset.CoreV1().Pods(namespace).Get(pod, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			fmt.Printf("Pod %s in namespace %s not found\n", pod, namespace)
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			fmt.Printf("Error getting pod %s in namespace %s: %v\n",
				pod, namespace, statusError.ErrStatus.Message)
		} else if err != nil {
			panic(err.Error())
		} else {
			fmt.Printf("Found pod %s in namespace %s\n", pod, namespace)
		}

		time.Sleep(10 * time.Second)
	}

}
