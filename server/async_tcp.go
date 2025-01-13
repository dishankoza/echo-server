package server

import (
	"log"
	"net"
	"syscall"

	"github.com/dishankoza/echo-server/config"
	"github.com/dishankoza/echo-server/core"
)

var con_clients int = 0

// RunAsyncTCPServer starts an asynchronous TCP server using EPOLL
func RunAsyncServer() error {
	log.Println("Starting an asynchronous TCP server on", config.Host, config.Port)
	max_clients := 20000

	// Create EPOLL Event Objects to hold events
	var events []syscall.EpollEvent = make([]syscall.EpollEvent, max_clients)

	// Create a socket
	serverFD, err := syscall.Socket(syscall.AF_INET, syscall.O_NONBLOCK|syscall.SOCK_STREAM, 0)
	if err != nil {
		return err
	}
	defer syscall.Close(serverFD)

	// Set the Socket to operate in non-blocking mode
	if err = syscall.SetNonblock(serverFD, true); err != nil {
		return err
	}

	// Bind the IP and the port
	ip4 := net.ParseIP(config.Host)
	if err = syscall.Bind(serverFD, &syscall.SockaddrInet4{
		Port: config.Port,
		Addr: [4]byte{ip4[0], ip4[1], ip4[2], ip4[3]},
	}); err != nil {
		return err
	}

	// Start listening on the socket
	if err = syscall.Listen(serverFD, max_clients); err != nil {
		return err
	}

	// Create EPOLL instance
	epollFD, err := syscall.EpollCreate1(0)
	if err != nil {
		log.Fatal(err)
	}
	defer syscall.Close(epollFD)

	// Add the server socket to the EPOLL instance to listen for new connections
	socketServerEvent := syscall.EpollEvent{
		Events: syscall.EPOLLIN, // Monitor for read events (incoming connections)
		Fd:     int32(serverFD),
	}
	if err = syscall.EpollCtl(epollFD, syscall.EPOLL_CTL_ADD, serverFD, &socketServerEvent); err != nil {
		return err
	}

	for {
		// Wait for events (blocking call)
		nevents, e := syscall.EpollWait(epollFD, events[:], -1)
		if e != nil {
			continue
		}

		for i := 0; i < nevents; i++ {
			// If the server socket is ready for an IO (new connection)
			if int(events[i].Fd) == serverFD {
				// Accept the incoming client connection
				fd, _, err := syscall.Accept(serverFD)
				if err != nil {
					log.Println("Error accepting connection:", err)
					continue
				}

				// Increase the number of concurrent clients
				con_clients++
				syscall.SetNonblock(fd, true) // Set client socket to non-blocking

				// Add this new client socket to the EPOLL instance to monitor it
				socketClientEvent := syscall.EpollEvent{
					Events: syscall.EPOLLIN, // Monitor for read events (client data)
					Fd:     int32(fd),
				}
				if err := syscall.EpollCtl(epollFD, syscall.EPOLL_CTL_ADD, fd, &socketClientEvent); err != nil {
					log.Fatal("Failed to add client socket to epoll:", err)
					continue
				}

			} else {
				// Handle data from an existing client
				comm := core.FDComm{Fd: int(events[i].Fd)}
				cmd, err := readCommandFD(comm)
				if err != nil {
					// Handle read error (client disconnected)
					log.Println("Error reading from client:", err)
					syscall.Close(int(events[i].Fd))
					con_clients--
					continue
				}

				// Respond to the client with the command
				respondFD(cmd, comm)
			}
		}
	}
}

// readCommand reads the data from the client and returns the command
func readCommandFD(comm core.FDComm) (string, error) {
	// Define a buffer to store the command
	buffer := make([]byte, 512)
	n, err := comm.Read(buffer)
	if err != nil {
		return "", err
	}
	return string(buffer[:n]), nil
}

// respond sends the response back to the client
func respondFD(cmd string, comm core.FDComm) {
	// Send the response back to the client
	_, err := comm.Write([]byte("Echo: " + cmd))
	if err != nil {
		log.Println("Error writing to client:", err)
		syscall.Close(comm.Fd)
	}
}
