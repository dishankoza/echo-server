package main

import (
	"flag"
	"log"

	"github.com/dishankoza/echo-server/config"
	"github.com/dishankoza/echo-server/server"
)

func setupFlags() {
	flag.StringVar(&config.Host, "host", "0.0.0.0", "host for the echo server")
	flag.IntVar(&config.Port, "port", 7379, "port for the server")
	flag.Parse()
}

func main() {
	setupFlags()
	log.Println("Echo Server up and listening")
	go server.RunASyncServerWRoutine()
	go server.RunSyncServer()
	server.RunASyncServerWNetpoll()
}
