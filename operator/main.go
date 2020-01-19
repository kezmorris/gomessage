package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

// Operator object to listen and serve
type Operator struct {
	listener    net.Listener
	logger      *Logger
	kubeclient  *KubeClient
	conferences []*Conference
}

// Conference is a single call conference.
type Conference struct {
	port int
	name string
}

// NewOperator creates a new Operator
func NewOperator(port int) (*Operator, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return nil, fmt.Errorf("Failed to during net.Listen: %v", err)

	}

	logger := NewLogger()
	logger.Run()

	var client *KubeClient

	if os.Getenv("KUBECONFIG") == "IN_CLUSTER" {
		client, err = NewInClusterKubeClient(&logger)
	} else {
		client, err = NewKubeClient(&logger)
	}
	if err != nil {
		return nil, fmt.Errorf("Error creating kubernetes client: %v", err)
	}
	client.Run()

	srv := &Operator{
		ln,
		&logger,
		client,
	}
	return srv, nil
}

// Serve serves as the Operator
func (srv *Operator) Serve() error {
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

// Flow:
// Create call pod, giving it the token to be used for login
// Send acknowledgement with port to connect to
// Client will wait until pod is ready

func (srv *Operator) createCallSession(conn net.Conn) {
	connID := conn.RemoteAddr().(*net.TCPAddr).Port
	srv.logger.Log(fmt.Sprintf("%v wants to connect", connID))
	sessionPort, token := srv.getSessionPort()
	srv.logger.Log(fmt.Sprintf("Directing %v to port %v", connID, sessionPort))

	newConference := srv.createConference(sessionPort, token)

	srv.conferences.append(srv.conferences, newConference)

	response := fmt.Sprintf("CONN %v %v %v", sessionPort, token, newConference)
	conn.Write([]byte(response + "\n"))
}

func (srv *Operator) createConference(port int, token int) *Conference {
	newConf := &Conference{}
	return newConf
}

func (srv *Operator) getSessionPort() (int, int) {
	port := srv.kubeclient.RequestPort()
	return port, 123456
}

func main() {
	operator, err := NewOperator(8001)
	if err != nil {
		log.Panicf("Error during creation of operator: %v", err)
	}
	operator.logger.Log("Starting ")

	operator.Serve()

}
