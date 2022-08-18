package main

import (
	"fmt"
	"github.com/Rock-liyi/p2pdb-pubsub/chat/internal/server"
	"net/http"
)

func main() {
	serve := server.NewServer()
	mux := http.NewServeMux()

	go serve.Run()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Welcome"))
	})

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		server.ServeWs(serve, w, r)
	})

	fmt.Println("http://127.0.0.1:8080, running...")
	_ = http.ListenAndServe(":8080", mux)
}
