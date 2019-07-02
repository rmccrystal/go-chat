package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type options struct {
	serverIp   string
	serverPort int
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Enter IP Address: ")
	ip, _ := reader.ReadString('\n')

	// Remove all newlines in ip
	ip = strings.ReplaceAll(ip, "\n", "")

	port := 0
	for {
		fmt.Printf("Enter port: ")
		portInput, _ := reader.ReadString('\n')

		// Remove all newlines in the input
		portInput = strings.ReplaceAll(portInput, "\n", "")

		// Convert the port entered by the user into an int
		p, err := strconv.Atoi(portInput)

		// Doing this because of scoping challenges
		// The port variable contains the actual port and this sends the p variable to that port variable
		port = p
		if err != nil || port > 65535 || port < 1 {
			fmt.Printf("%s is not a valid port\n", portInput)
			continue
		}
		break
	}

	chat(options{ip, port})
}

func chat(opts options) {
	conn, err := net.Dial("tcp",
		fmt.Sprintf("%s:%d", opts.serverIp, opts.serverPort))

	if err != nil {
		fmt.Println(err)
	}

	// Make a new thread to write data to stdout
	go func(conn net.Conn) {
		for {
			data, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				fmt.Println(err)
			}
			fmt.Print(data)
		}
	}(conn)

	writer := bufio.NewWriter(conn)
	_, err = writer.ReadFrom(os.Stdin)

	if err != nil {
		fmt.Println(err)
	}
}
