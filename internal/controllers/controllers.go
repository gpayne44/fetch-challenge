package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gpayne44/fetch-challenge/internal/repositories"
)

type controller struct {
	repository repositories.ReceiptsRepository
	logger     log.Logger
}

func New(repository repositories.ReceiptsRepository) *controller {
	return &controller{
		repository: repository,
		logger:     *log.Default(),
	}
}

func (c *controller) Register(router *mux.Router) {
	router.HandleFunc("/receipts/process", c.ProcessReceipt()).Methods(http.MethodPost)
	router.HandleFunc("/receipts/{id}/points", c.GetReceiptPoints()).Methods(http.MethodGet)
}
