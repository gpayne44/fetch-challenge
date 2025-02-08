package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gpayne44/fetch-challenge/internal/entities"
	"github.com/gpayne44/fetch-challenge/internal/process"
	"github.com/gpayne44/fetch-challenge/internal/repositories"
)

func (c *controller) ProcessReceipt() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			c.logger.Println("error reading request body")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("error reading request body: %s", err.Error())))
			return
		}

		var receipt entities.Receipt
		err = json.Unmarshal(b, &receipt)
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("could not unmarshal request: %s", err.Error())))
			return
		}

		valid := receipt.Validate()
		if !valid {
			c.logger.Println("The receipt is invalid.")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("The receipt is invalid."))
			return
		}

		pointTotal, processErrors := process.CalculatePoints(receipt)
		if len(processErrors) != 0 {
			c.logger.Printf("error calculating point total: %v", processErrors)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("error calculating point total: %v", processErrors)))
			return
		}

		record := entities.ReceiptRecord{Points: pointTotal, Receipt: receipt}
		newID, err := c.repository.StoreReceipt(record)
		if err != nil {
			c.logger.Println("error storing receipt")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("error storing receipt: %s", err.Error())))
			return
		}

		resBytes, err := json.Marshal(entities.ProcessResponse{ID: newID})
		if err != nil {
			c.logger.Printf("receipt processed, could not marshal new ID: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("receipt processed, could not marshal new ID: %s", err.Error())))
			return
		}
		w.Write(resBytes)
	}
}

func (c *controller) GetReceiptPoints() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		if id == "" {
			c.logger.Println("empty ID in request")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("empty ID in request"))
			return
		}

		record, err := c.repository.GetReceipt(id)
		if err != nil {
			if err == repositories.ErrNotFound {
				c.logger.Printf("record not found for id %s", id)
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("No receipt found for that ID."))
				return
			} else if err != repositories.ErrNotFound {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("error reading record for id %s: %s", id, err.Error())))
				return
			}
		}

		resBytes, err := json.Marshal(entities.PointsResponse{Points: record.Points})
		if err != nil {
			c.logger.Printf("could not marhsal response: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("could not marhsal response: %s", err.Error())))
			return
		}
		w.Write(resBytes)
	}
}
