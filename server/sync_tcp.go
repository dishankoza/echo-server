package server

import (
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/dishankoza/echo-server/config"
)

func readCommand(c net.Conn) (string, error) {
	//TODO making the read shot compatible with size > 512

	var buf []byte = make([]byte, 512)
	n, err := c.Read(buf[:])
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}

func respond(cmd string, c net.Conn) error {
	if _, err := c.Write([]byte(cmd)); err != nil {
		return err
	}
	return nil
}

func RunSyncServer() {
	log.Println("Starting a synchronous TCP server on", config.Host, config.Port)

	var con_client int = 0

	lsnr, err := net.Listen("tcp", config.Host+":"+strconv.Itoa((config.Port)))

	if err != nil {
		panic(err)
	}

	for {
		c, err := lsnr.Accept()
		if err != nil {
			panic(err)
		}

		con_client += 1
		log.Println("Client connected with address", c.RemoteAddr(), "concurrent clients", con_client)

		for {
			cmd, err := readCommand(c)
			if err != nil {
				c.Close()
				con_client -= 1
				log.Println("Client Disconnected", c.RemoteAddr(), "concurrent clients", con_client)
				if err == io.EOF {
					break
				}
				log.Println("err", err)
			}
			cmds := strings.Split(cmd, "\n")
			for _, cmd := range cmds {
				if cmd != "" {
					log.Println("Command", cmd)
				}
			}

			if err = respond(cmd, c); err != nil {
				log.Println("err write:", err)
			}
		}
	}
}
