package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Operator object to listen and serve
type Operator struct {
	listener net.Listener
	logger   Logger
}

// NewOperator creates a new Operator
func NewOperator(port int) (*Operator, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return nil, fmt.Errorf("Failed to during net.Listen: %v", err)

	}
	logger := NewLogger()
	logger.Run()
	srv := &Operator{
		ln,
		logger,
	}
	return srv, nil
}

// Serve serves as the Operator
func (srv *Operator) Serve() error {
	var cs *KubeClient
	var err error

	if os.Getenv("KUBECONFIG") == "IN_CLUSTER" {
		cs, err = NewInClusterKubeClient()
	} else {
		cs, err = NewKubeClient()
	}
	if err != nil {
		return fmt.Errorf("Error creating kubernetes client: %v", err)
	}

	go srv.doSomethingWithClient(cs)

	for {

		conn, err := srv.listener.Accept()
		if err != nil {
			return fmt.Errorf("Error accepting connection: %v", err)
		}
		go srv.handleConnection(conn)
	}
}

func (srv *Operator) handleConnection(conn net.Conn) error {
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				srv.logger.Log(fmt.Sprintf("Error whilst connecting to %v: %v", conn.RemoteAddr().(*net.TCPAddr).Port, err))
				break
			}
		}

		message = strings.TrimSpace(message)
		if message == "01" {
			srv.createCallSession(conn)
		} else {
			// Ignore
		}
	}
	return nil
}

func (srv *Operator) createCallSession(conn net.Conn) {
	connID := conn.RemoteAddr().(*net.TCPAddr).Port
	srv.logger.Log(fmt.Sprintf("%v wants to connect", connID))
	sessionPort, token := getSessionPort()
	srv.logger.Log(fmt.Sprintf("Directing %v to port %v", connID, sessionPort))
	response := fmt.Sprintf("CONN %v %v", sessionPort, token)
	conn.Write([]byte(response + "\n"))
}

func getSessionPort() (int, int) {
	return 30001, 123456
}

func (srv *Operator) doSomethingWithClient(client *KubeClient) {
	var namespace string

	if os.Getenv("NAMESPACE") != "" {
		namespace = os.Getenv("NAMESPACE")
	} else {
		namespace = "" // Either the deployment is misconfigured or not running in cluster
	}
	for {
		pods, err := client.clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{})
		if err != nil {
			srv.logger.Log(fmt.Sprintf("Error calling clientset: %v", err))
		}
		srv.logger.Log(fmt.Sprintf("There are %d pods in this namespace\n", len(pods.Items)))
		time.Sleep(20 * time.Second)
	}
}

func main() {
	operator, err := NewOperator(8001)
	if err != nil {
		log.Panicf("Error during creation of operator: %v", err)
	}
	operator.logger.Log("Starting ")

	operator.Serve()

}
