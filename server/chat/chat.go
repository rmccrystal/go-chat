package chat

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

const (
	MinUsernameLength int = 3
	MaxUsernameLength int = 16
)

type user struct {
	username string
	admin    bool
}

var users = make(map[net.Conn]user)

func HandleNewConnection(conn net.Conn) error {
	return getUsername(conn)
}

func HandleEndConnection(conn net.Conn) {
	// The exclude value is user{} to specify we don't want to exclude anyone
	err := broadcast(fmt.Sprintf("%s left the chat\n", users[conn].username), user{})
	if err != nil {
		fmt.Println(err)
	}
	delete(users, conn)
}

func HandleMessage(conn net.Conn, message string) error {
	// If the message is too short, don't do anything
	if len(message) <= 1 {
		return nil
	}

	// Broadcast message to everyone except the sender
	err := broadcast(fmt.Sprintf("%s: %s", users[conn].username, message), users[conn])
	if err != nil {
		return err
	}
	return nil
}

// If exclude is not nil, the message will not be sent to that specific user
func broadcast(msg string, exclude user) error {
	for conn, usr := range users {
		// If we hit the user we need to exclude then continue the loop without sending data to them
		if usr == exclude {
			continue
		}

		// Write message
		_, err := conn.Write([]byte(msg))
		if err != nil {
			return err
		}
	}

	// No errors
	return nil
}

// Get username from connection
func getUsername(conn net.Conn) error {
	// Loop until we get a valid username
	for {
		conn.Write([]byte("Enter your username: "))

		// Read result
		username, err := bufio.NewReader(conn).ReadString('\n')

		if err != nil {
			return err
		}

		// Remove all newlines in username
		username = strings.ReplaceAll(username, "\n", "")

		// Check if the username is too short
		if len(username) > MaxUsernameLength {
			conn.Write([]byte(fmt.Sprintf(
				"The username %s is longer than %d characters\n",
				username, MaxUsernameLength)))
			continue
		}

		// Check if the username is too long
		if len(username) < MinUsernameLength {
			conn.Write([]byte(fmt.Sprintf(
				"The username %s is shorter than %d characters\n",
				username, MinUsernameLength)))
			continue
		}

		// Check if the username is already taken
		isUsernameTaken := false // true if the username is already taken
		for _, user := range users {
			if username == user.username {
				conn.Write([]byte(fmt.Sprintf(
					"The username %s is already taken\n",
					username)))
				isUsernameTaken = true
			}
		}
		// If the username is taken continue with the for loop
		if isUsernameTaken {
			continue
		}

		// Add new user
		users[conn] = user{username: username, admin: false}

		// Tell everyone someone joined
		// The user{} in the second parameter of broadcast() specifies that we don't want to exclude anyone from this message
		err = broadcast(fmt.Sprintf("%s joined the chat\n", username), user{})
		if err != nil {
			return err
		}

		return nil
	}
}
