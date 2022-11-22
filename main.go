package main

import (
	"log"

	"github.com/seagalputra/talkbox/api"
)

func main() {
	if err := api.StartServer(); err != nil {
		log.Panicf("Failed to start the server: %v", err)
	}
}
