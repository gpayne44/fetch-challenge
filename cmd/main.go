package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gpayne44/fetch-challenge/internal/controllers"
	"github.com/gpayne44/fetch-challenge/internal/repositories"
)

func main() {
	m := repositories.New()
	c := controllers.New(m)

	r := mux.NewRouter()
	c.Register(r)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
	}

	srv.ListenAndServe()
}
