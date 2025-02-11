package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gpayne44/fetch-challenge/internal/repositories"
)

var (
	errFmtReadingRequest    = "error reading request body: %s"
	errFmtUnmarshalRequest  = "could not unmarshal request: %s"
	errFmtCalculatePoints   = "error calculating point total: %v"
	errFmtStoreReceipt      = "error storing receipt: %s"
	errFmtMarshalIDResponse = "receipt processed, could not marshal new ID: %s"
	errFmtReceiptReadError  = "error reading record for id %s: %s"
	errFmtMarshalResponse   = "could not marhsal response: %s"
	errFmtInvalidReceiptID  = "could not parse id param %s: %s"

	errMsgInvalidReceipt = "The receipt is invalid."
	errEmptyID           = "empty ID in request path"
	errNoReceiptFound    = "No receipt found for that ID."
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
	router.HandleFunc("/receipts/{id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}/points", c.GetReceiptPoints()).Methods(http.MethodGet)
}
