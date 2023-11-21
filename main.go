package main

import (
	"log"
	"os"

	"github.com/raja-dettex/goatQ_consumer/api"
	"github.com/raja-dettex/goatQ_consumer/server"
)

func main() {
	addr := os.Getenv("ADDR")
	listenAddr := os.Getenv("LISTEN_ADDR")
	channel := make(chan []byte)
	opts := server.ConsumerOpts{Addr: addr}
	consumer := server.NewGoatQConsumer(opts, channel)
	go consumer.Subscribe()
	serverOpts := api.ServerOpts{ListenAddr: listenAddr}
	server := api.NewAPIServer(serverOpts, consumer)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
