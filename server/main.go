package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

// Server object to listen and serve
type Server struct {
	listener net.Listener
}

// NewServer creates a new server
func NewServer(port int) (*Server, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return nil, fmt.Errorf("Failed to during net.Listen: %v", err)

	}

	srv := &Server{
		ln,
	}
	return srv, nil
}

// Serve serves as the server
func (srv *Server) Serve() error {
	messages := NewLogger()
	messages.Run()
	for {

		conn, err := srv.listener.Accept()
		if err != nil {
			return fmt.Errorf("Error accepting connection: %v", err)
		}

		log.Printf("Setup listener successfully. New client: %v", conn.RemoteAddr().(*net.TCPAddr).Port)

		go srv.handleConnection(conn, messages.logChan)
	}
}
func (srv *Server) handleConnection(conn net.Conn, logChan chan string) error {
	log.Printf("Starting to listen..")

	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return fmt.Errorf("Error whilst listening, %v", err)
			}
		}
		logChan <- strings.TrimSpace(message)

	}
	log.Printf("Connection closed")
	return nil
}

func main() {
	log.Printf("Starting listener server")

	server, err := NewServer(8001)
	if err != nil {
		log.Panicf("Failed to during ln. Accept: %v", err)
	}

	server.Serve()

}
