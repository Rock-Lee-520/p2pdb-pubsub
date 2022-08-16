package main

import (
	"fmt"
	"github.com/Rock-liyi/p2pdb-pubsub/chat/internal"
	"net/http"
)

func main() {
	hub := internal.NewHub()
	server := http.NewServeMux()

	go hub.Run()

	server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Welcome"))
	})

	server.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		internal.ServeWs(hub, w, r)
	})

	fmt.Println("http://127.0.0.1:8080, running...")
	_ = http.ListenAndServe(":8080", server)
}
