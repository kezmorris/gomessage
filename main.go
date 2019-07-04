package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	log.Printf("Starting listener server")
	conn, err := setup()
	if err != nil {
		log.Panicf("Error while setting up: %v", err)
	}
	log.Printf("Setup listener successfully.")
	err = listenLoop(conn)
	if err != nil {
		log.Panicf("Error while listening: %v", err)
	}
}

func setup() (net.Conn, error) {
	ln, err := net.Listen("tcp", ":8001")
	if err != nil {
		return nil, fmt.Errorf("Failed to during net.Listen: %v", err)

	}
	conn, err := ln.Accept()
	if err != nil {
		return nil, fmt.Errorf("Failed to during ln.Accept: %v", err)
	}
	return conn, nil
}

func listenLoop(conn net.Conn) error {
	log.Printf("Starting to listen..")

	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			return fmt.Errorf("Error whilst listening, ", err)
		}
		log.Print("Message Received:", string(message))

		newmessage := strings.ToUpper(message)
		conn.Write([]byte(newmessage + "\n"))

	}
}
