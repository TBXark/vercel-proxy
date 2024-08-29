package main

import (
	"github.com/tbxark/vercel-proxy/api"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", api.Handler)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
		return
	}
}
