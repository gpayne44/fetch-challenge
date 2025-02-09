package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/gpayne44/fetch-challenge/internal/controllers"
	"github.com/gpayne44/fetch-challenge/internal/repositories"
)

func main() {
	m := repositories.New()
	c := controllers.New(m)

	r := mux.NewRouter()
	c.Register(r)

	var port string
	flag.StringVar(&port, "port", "8000", "localhost port")
	flag.Parse()
	addr := fmt.Sprintf("127.0.0.1:%s", port)

	srv := &http.Server{
		Handler: r,
		Addr:    addr,
	}

	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Println("Server shutting down...")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		srv.Close()
		log.Fatalf("Error shutting down server: %v", err)
	}
	log.Println("Server shutdown complete.")
}
