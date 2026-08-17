package main

import (
	"log"
	"net/http"
	"os"

	"github.com/zhangkui/go-feature-rollout/internal/httpapi"
	"github.com/zhangkui/go-feature-rollout/internal/service"
)

func main() {
	address := os.Getenv("HTTP_ADDR")
	if address == "" {
		address = ":8080"
	}
	server := &http.Server{
		Addr:    address,
		Handler: httpapi.NewHandler(service.New()),
	}
	log.Printf("feature rollout service listening on %s", address)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
