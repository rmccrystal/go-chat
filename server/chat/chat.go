package chat

import (
	"../logger"
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

const (
	MinUsernameLength int = 3
	MaxUsernameLength int = 16
)

type user struct {
	username string
	admin    bool
}

var usersMutex = sync.Mutex{}
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

	// usersMutex.Lock and Unlock are for thread safety
	usersMutex.Lock()
	delete(users, conn)
	usersMutex.Unlock()
}

func HandleMessage(conn net.Conn, message string) error {
	// If the message is too short, don't do anything
	if len(message) <= 1 {
		return nil
	}

	// Thread safety
	// Get the data that should be sent to the broadcast function while the mutex is locked
	usersMutex.Lock()
	broadcastMsg := fmt.Sprintf("%s: %s", users[conn].username, message)
	excludeUser := users[conn]
	usersMutex.Unlock()
	//

	// Broadcast message to everyone except the sender
	err := broadcast(broadcastMsg, excludeUser)
	if err != nil {
		return err
	}
	return nil
}

// If exclude is not nil, the message will not be sent to that specific user
func broadcast(msg string, exclude user) error {
	// Log the message
	logger.LogPrefixf(msg, "CHAT")

	// thread safety
	usersMutex.Lock()
	defer usersMutex.Unlock()

	for conn, usr := range users {
		// If we hit the user we need to exclude then continue the loop without sending data to them
		if usr == exclude {
			continue
		}

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

		// Lock users mutex for thread safety
		usersMutex.Lock()

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
			// If the user needs to select another username, unlock the mutex
			usersMutex.Unlock()
			continue
		}

		// Add new user
		users[conn] = user{username: username, admin: false}

		// Now that we've added the user to the list, we can unlock the mutex
		usersMutex.Unlock()

		// Tell everyone someone joined
		// The user{} in the second parameter of broadcast() specifies that we don't want to exclude anyone from this message
		err = broadcast(fmt.Sprintf("%s joined the chat\n", username), user{})
		if err != nil {
			return err
		}

		return nil
	}
}
