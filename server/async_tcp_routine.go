package server

import (
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/dishankoza/echo-server/config"
)

type ConnWrapper struct {
	net.Conn
}

func (c *ConnWrapper) readCommand() (string, error) {
	var buf []byte
	tmp := make([]byte, 512)
	for {
		n, err := c.Read(tmp)
		if err != nil {
			log.Print(err == io.EOF)
			if err == io.EOF {
				return "", err
			}
			return "", err
		}
		buf = append(buf, tmp[:n]...)
		if n < 512 {
			break
		}
	}
	return string(buf), nil
}

func RunASyncServerWRoutine() {
	log.Println("Starting a synchronous TCP server on", config.Host, config.Port)

	var conClient int
	var mu sync.Mutex // Mutex to protect conClient
	var md sync.Mutex

	lsnr, err := net.Listen("tcp", config.Host+":"+strconv.Itoa((config.Port)))

	if err != nil {
		panic(err)
	}

	for {
		c, err := lsnr.Accept()
		if err != nil {
			panic(err)
		}

		mu.Lock() // Acquire the mutex before modifying conClient
		conClient++
		log.Println("Client connected with address", c.RemoteAddr(), "concurrent clients", conClient)
		mu.Unlock() // Release the mutex after modifying conClient

		go func(c net.Conn) {
			defer c.Close()
			defer func() {
				md.Lock()
				conClient--
				log.Println("Client disconnected", c.RemoteAddr(), "concurrent clients", conClient)
				md.Unlock()
			}()

			wrappedConn := &ConnWrapper{c}
			for {
				cmd, err := wrappedConn.readCommand()
				if err != nil {
					if err == io.EOF {
						return
					}
					log.Println("Failed to read command:", err)
					return
				}

				cmds := strings.Split(cmd, "\n")
				for _, cmd := range cmds {
					if cmd != "" {
						log.Println("Command", cmd)
					}
				}

				if err = respond(cmd, c); err != nil {
					log.Println("Failed to respond:", err)
					return
				}
			}
		}(c)
	}
}
