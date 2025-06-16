package main

import (
	"demo/server"
	"log"
)

func main() {
	if err := server.StartServer(":8090"); err != nil {
		log.Fatal(err)
	}
}
