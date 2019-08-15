package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

// Operator object to listen and serve
type Operator struct {
	listener net.Listener
}

// NewOperator creates a new Operator
func NewOperator(port int) (*Operator, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return nil, fmt.Errorf("Failed to during net.Listen: %v", err)

	}

	srv := &Operator{
		ln,
	}
	return srv, nil
}

// Serve serves as the Operator
func (srv *Operator) Serve() error {
	logman := NewLogger()
	logman.Run()
	for {

		conn, err := srv.listener.Accept()
		if err != nil {
			return fmt.Errorf("Error accepting connection: %v", err)
		}
		go srv.handleConnection(conn, logman)
	}
}

func (srv *Operator) handleConnection(conn net.Conn, logman Logger) error {
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				logman.Log(fmt.Sprintf("Error whilst connecting to %v: %v", conn.RemoteAddr().(*net.TCPAddr).Port, err))
				break
			}
		}

		message = strings.TrimSpace(message)
		if message == "01" {
			createCallSession(conn, logman)
		} else {
			// Ignore
		}
	}
	return nil
}

func createCallSession(conn net.Conn, logman Logger) {
	connID := conn.RemoteAddr().(*net.TCPAddr).Port
	logman.Log(fmt.Sprintf("%v wants to connect", connID))
	sessionPort, token := getSessionPort()
	logman.Log(fmt.Sprintf("Directing %v to port %v", connID, sessionPort))
	response := fmt.Sprintf("CONN %v %v", sessionPort, token)
	conn.Write([]byte(response + "\n"))
}

func getSessionPort() (int, int) {
	return 30001, 123456
}

func main() {
	log.Printf("Starting operator")

	operator, err := NewOperator(8001)
	if err != nil {
		log.Panicf("Failed to during ln. Accept: %v", err)
	}

	operator.Serve()

}
