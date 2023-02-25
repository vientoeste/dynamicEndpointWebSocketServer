package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// map of dynamic endpoints to handler functions
var endpoints = make(map[string]func(*websocket.Conn))

func main() {
	// register endpoint handlers
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ws", wsHandler)

	fmt.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the home page!")
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// HTTP > WebSocket upgrader
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		// Print message to console
		fmt.Println(string(p))

		// Route message to appropriate handler function
		switch string(p) {
		case "register":
			registerEndpoint(conn)
		default:
			handleEndpoint(string(p), conn)
		}
	}
}

func registerEndpoint(conn *websocket.Conn) {
	// read new endpoint name from client
	_, p, err := conn.ReadMessage()

	if err != nil {
		fmt.Println(err)
		return
	}

	// Create new handler function for new endpoint
	endpointName := string(p)
	endpointFunc := func(conn *websocket.Conn) {
		for {
			// Read message from client
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				return
			}

			// Print message to console
			fmt.Println(string(p))

			// Send message back to client
			err = conn.WriteMessage(messageType, p)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}

	// add new endpoint to the map
	endpoints[endpointName] = endpointFunc
}

func handleEndpoint(endpoint string, conn *websocket.Conn) {
	// get handler function for the endpoint from the map
	endpointFunc, ok := endpoints[endpoint]
	if !ok {
		fmt.Println("Invalid endpoint")
		return
	}

	// call handler function for endpoint
	endpointFunc(conn)
}
