package main

import (
	"gitlab.com/mfcekirdek/in-memory-store/configs"
)

func main() {
	config := configs.NewConfig()
	config.Print()
}
