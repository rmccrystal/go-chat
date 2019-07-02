package main

import (
	"./chat"
	"bufio"
	"fmt"
	"net"
)

// Config
const (
	ServerPort int = 8080
)

func main() {
	// Start listening
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", ServerPort))

	// Check for errors
	if err != nil {
		fmt.Printf("Error creating listener: %s\n", err.Error())
		return
	}

	// Close the listener when the server stops0.
	defer ln.Close()
	fmt.Printf("go-chat server listening on port %d\n", ServerPort)

	// Run loop that listens for connections
	listenerLoop(ln)
}

func listenerLoop(ln net.Listener) {
	for {
		// Accept connections
		conn, err := ln.Accept()

		if err != nil {
			fmt.Printf("error accepting connection from %s\n", conn.RemoteAddr().String())
		}

		// Create a new goroutine to handle the connection
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	// Automatically close the connection when we're done with it
	defer conn.Close()

	fmt.Printf("Handling connection from %s\n", conn.RemoteAddr().String())

	defer chat.HandleEndConnection(conn)

	err := chat.HandleNewConnection(conn)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		// Read string until we encounter a newline
		data, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		err = chat.HandleMessage(conn, data)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
