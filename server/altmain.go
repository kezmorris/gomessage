package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

// Basic flow: The client will signal the main port (8001) with it's name, asking for a session.
// The server will respond with a port for the client to connect to.
// If the users group doesn't have a session, it will make one. If it does, it will simply send
// the port for that session. Multiple messages may be sent.

// The hard part here, is making this all work statelessly. I think that perhaps the initial
// server will be a lone pod, with a connected (or local DB), and each session will be stateful,
// and started up by the server.

//

func main() {
	log.Printf("Starting listener server")
	for {
		conn, err := setup(8081)
		if err != nil {
			log.Panicf("Error while setting up: %v", err)
		}

		log.Printf("Setup listener successfully.")
		result, err := doSession(conn)
		if err != nil {
			log.Panicf("Session closed: %v", err)
		}
		log.Printf(result)
	}
}

func setup(port int) (net.Conn, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return nil, fmt.Errorf("Failed to during net.Listen: %v", err)

	}
	conn, err := ln.Accept()
	if err != nil {
		return nil, fmt.Errorf("Failed to during ln. Accept: %v", err)
	}
	return conn, nil
}

func doSession(conn net.Conn) (string, error) {
	log.Printf("Starting to listen..")

	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return "", fmt.Errorf("Error whilst listening, %v", err)
			}
		}
		log.Print("Message Received:", string(message))

		newmessage := strings.ToUpper(message)
		conn.Write([]byte(newmessage + "\n"))

	}
	return "Completed session.", nil
}
