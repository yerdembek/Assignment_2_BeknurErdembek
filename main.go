package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yerdembek/Assignment_2_BeknurErdembek/internal/server"
	"github.com/yerdembek/Assignment_2_BeknurErdembek/internal/store"
)

func main() {
	s := store.NewStore[string, string]()
	srv := server.New(":8080", s)

	go func() {
		log.Println("server starting on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)

	<-stopCh
	log.Println("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	} else {
		log.Println("server stopped gracefully")
	}
}
