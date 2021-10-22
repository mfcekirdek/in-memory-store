// Main package loads the configuration
package main

import (
	"log"

	"gitlab.com/mfcekirdek/in-memory-store/configs"
	"gitlab.com/mfcekirdek/in-memory-store/internal"
)

func main() {
	config := configs.NewConfig()
	config.Print()

	server := internal.NewServer(config)
	err := server.Start()
	if err != nil {
		log.Fatal(err)
	}
}
