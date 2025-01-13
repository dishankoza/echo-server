package server

import (
	"context"
	"io"
	"log"
	"sync"

	"github.com/cloudwego/netpoll"
)

var (
	conClient int
	mu        sync.Mutex // Mutex to protect concurrent access to conClient
)

// fetchCommand reads data from the connection
func fetchCommand(c netpoll.Connection) (string, error) {
	var buf []byte = make([]byte, 512)
	n, err := c.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}

// handleConnection handles each client connection
func handleConnection(c netpoll.Connection) {
	defer c.Close()

	// Increment the client count
	mu.Lock()
	conClient++
	log.Println("Client connected with address", c.RemoteAddr(), "concurrent clients", conClient)
	mu.Unlock()

	for {
		// Fetch command from the client
		cmd, err := fetchCommand(c)
		if err != nil {
			if err != io.EOF {
				log.Println("read error:", err)
			}
			break
		}
		log.Println("Received command from "+c.RemoteAddr().String()+":", cmd)

		// Respond with the same command
		if err = respond(cmd, c); err != nil {
			log.Println("err write:", err)
			break
		}
	}

	// Decrement the client count when done
	mu.Lock()
	conClient--
	log.Println("Client disconnected with address", c.RemoteAddr(), "concurrent clients", conClient)
	mu.Unlock()
}

// RunASyncServerWNetpoll starts the asynchronous server
func RunASyncServerWNetpoll() {
	log.Println("Starting asynchronous TCP server on port 8080")

	// Create the listener
	listener, err := netpoll.CreateListener("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	// Create the event loop
	eventLoop, err := netpoll.NewEventLoop(func(ctx context.Context, connection netpoll.Connection) error {
		// Handle each connection
		handleConnection(connection)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// Start the event loop to serve the listener
	if err := eventLoop.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
