package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gpayne44/fetch-challenge/internal/entities"
	"github.com/gpayne44/fetch-challenge/internal/process"
	"github.com/gpayne44/fetch-challenge/internal/repositories"
)

func (c *controller) ProcessReceipt() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			c.logger.Printf(errFmtReadingRequest, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(errFmtReadingRequest, err.Error())))
			return
		}

		var receipt entities.Receipt
		err = json.Unmarshal(b, &receipt)
		if err != nil {
			c.logger.Printf(errFmtUnmarshalRequest, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(errFmtUnmarshalRequest, err.Error())))
			return
		}

		valid := receipt.Validate()
		if !valid {
			c.logger.Println(errMsgInvalidReceipt)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errMsgInvalidReceipt))
			return
		}

		pointTotal, processErrors := process.CalculatePoints(receipt)
		if len(processErrors) != 0 {
			c.logger.Printf(errFmtCalculatePoints, processErrors)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(errFmtCalculatePoints, processErrors)))
			return
		}

		record := entities.ReceiptRecord{Points: pointTotal, Receipt: receipt}
		newID, err := c.repository.StoreReceipt(record)
		if err != nil {
			c.logger.Printf(errFmtStoreReceipt, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(errFmtStoreReceipt, err.Error())))
			return
		}

		resBytes, err := json.Marshal(entities.ProcessResponse{ID: newID})
		if err != nil {
			c.logger.Printf(errFmtMarshalIDResponse, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(errFmtMarshalIDResponse, err.Error())))
			return
		}
		w.Write(resBytes)
	}
}

func (c *controller) GetReceiptPoints() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := mux.Vars(r)["id"]
		if idParam == "" {
			c.logger.Println(errEmptyID)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errEmptyID))
			return
		}

		parsedID, err := uuid.Parse(idParam)
		if err != nil {
			c.logger.Println(errFmtInvalidReceiptID, idParam, err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf(errFmtInvalidReceiptID, idParam, err.Error())))
			return
		}

		record, err := c.repository.GetReceipt(parsedID)
		if err != nil {
			if err == repositories.ErrNotFound {
				c.logger.Println(errNoReceiptFound)
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(errNoReceiptFound))
				return
			} else if err != repositories.ErrNotFound {
				c.logger.Printf(errFmtReceiptReadError, parsedID.String(), err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf(errFmtReceiptReadError, parsedID.String(), err.Error())))
				return
			}
		}

		resBytes, err := json.Marshal(entities.PointsResponse{Points: record.Points})
		if err != nil {
			c.logger.Printf(errFmtMarshalResponse, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(errFmtMarshalResponse, err.Error())))
			return
		}
		w.Write(resBytes)
	}
}
