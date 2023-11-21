package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/raja-dettex/goatQ_consumer/server"
)

type ServerOpts struct {
	ListenAddr string
}

type APIServer struct {
	opts     ServerOpts
	consumer *server.GoatQConsumer
}

func NewAPIServer(opts ServerOpts, consumer *server.GoatQConsumer) *APIServer {
	return &APIServer{
		opts:     opts,
		consumer: consumer,
	}
}

func (server *APIServer) Start() error {
	server.registerHandlers()
	fmt.Printf("consumer listening on port %v\n", server.opts.ListenAddr)
	return http.ListenAndServe(server.opts.ListenAddr, nil)
}

func (server *APIServer) registerHandlers() {
	http.HandleFunc("/", server.handleConsume)
}

func (server *APIServer) handleConsume(w http.ResponseWriter, r *http.Request) {
	timeOutContext, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	for {
		select {
		case <-timeOutContext.Done():
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("error occurred"))
			return
		case msg := <-server.consumer.MessageChannel:
			w.WriteHeader(http.StatusOK)
			w.Write(msg)
			return
		}
	}

}
