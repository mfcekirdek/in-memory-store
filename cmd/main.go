package main

import (
	"gitlab.com/mfcekirdek/in-memory-store/configs"
	"gitlab.com/mfcekirdek/in-memory-store/internals"
	"log"
)

func main() {
	config := configs.NewConfig()
	config.Print()

	server := internals.NewServer(config)
	err := server.Start()
	if err != nil {
		log.Fatal(err)
	}
}
