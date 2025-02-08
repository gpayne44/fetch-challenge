package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gpayne44/fetch-challenge/internal/entities"
	"github.com/gpayne44/fetch-challenge/internal/repositories"
)

const (
	validReceipt = `
{
    "retailer": "Target",
    "purchaseDate": "2022-01-01",
    "purchaseTime": "13:01",
    "items": [
      {
        "shortDescription": "Mountain Dew 12PK",
        "price": "6.49"
      },
      {
        "shortDescription": "Emils Cheese Pizza",
        "price": "12.25"
      },
      {
        "shortDescription": "Knorr Creamy Chicken",
        "price": "1.26"
      },
      {
        "shortDescription": "Doritos Nacho Cheese",
        "price": "3.35"
      },
      {
        "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
        "price": "12.00"
      }
    ],
	"total": "35.35"
}
`
	invalidReceiptNoRetailer = `
{
    "retailer": "",
    "purchaseDate": "2022-01-01",
    "purchaseTime": "13:01",
    "items": [
      {
        "shortDescription": "Mountain Dew 12PK",
        "price": "6.49"
      },
      {
        "shortDescription": "Emils Cheese Pizza",
        "price": "12.25"
      },
      {
        "shortDescription": "Knorr Creamy Chicken",
        "price": "1.26"
      },
      {
        "shortDescription": "Doritos Nacho Cheese",
        "price": "3.35"
      },
      {
        "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
        "price": "12.00"
      }
    ],
	"total": "35.35"
}
`
	endpointProcess   = "/receipts/process"
	endpointGetPoints = "/receipts/%s/points"
)

func Test_ProcessReceipt(t *testing.T) {
	m := repositories.New()
	c := New(m)

	r := mux.NewRouter()
	c.Register(r)

	srv := httptest.NewServer(r)
	defer srv.Close()

	testCases := map[string]struct {
		input            string
		expStatusCode    int
		expectIDResponse bool
	}{
		"success with valid receipt": {
			input:            validReceipt,
			expStatusCode:    http.StatusOK,
			expectIDResponse: true,
		},
		"invalid receipt bad request": {
			input:         invalidReceiptNoRetailer,
			expStatusCode: http.StatusBadRequest,
		},
	}

	for caseName, tc := range testCases {
		t.Run(caseName, func(t *testing.T) {
			testData := map[string]interface{}{}
			err := json.Unmarshal([]byte(tc.input), &testData)
			if err != nil {
				t.Errorf("could not unmarshal test input: %s", err.Error())
				return
			}

			body, err := json.Marshal(testData)
			if err != nil {
				t.Errorf("could not marshal test input: %s", err.Error())
				return
			}

			res, err := http.Post(srv.URL+endpointProcess, "application/json", bytes.NewReader(body))
			if err != nil {
				t.Errorf("error sending test request: %s", err.Error())
				return
			}

			if res.StatusCode != tc.expStatusCode {
				t.Errorf("unexpected status code: got %d, want %d", res.StatusCode, tc.expStatusCode)
				return
			}
			if tc.expectIDResponse {
				b, err := io.ReadAll(res.Body)
				if err != nil {
					t.Errorf("error reading response body: %s", err.Error())
					return
				}

				var idRes entities.ProcessResponse
				err = json.Unmarshal(b, &idRes)
				if err != nil {
					t.Errorf("error unmarshal response body: %s", err.Error())
					return
				}
				if idRes.ID == "" {
					t.Error("empty ID in process response")
				}
			}
		})
	}
}

func Test_GetReceiptPoints(t *testing.T) {
	m := repositories.New()
	c := New(m)

	r := mux.NewRouter()
	c.Register(r)

	srv := httptest.NewServer(r)
	defer srv.Close()

	testData := map[string]interface{}{}
	err := json.Unmarshal([]byte(validReceipt), &testData)
	if err != nil {
		t.Errorf("could not unmarshal test input: %s", err.Error())
		return
	}

	testDataBytes, err := json.Marshal(testData)
	if err != nil {
		t.Fatal(err)
	}
	var receipt entities.Receipt
	err = json.Unmarshal(testDataBytes, &receipt)
	if err != nil {
		t.Fatal(err)
	}

	record := entities.ReceiptRecord{
		Receipt: receipt,
		Points:  28,
	}
	existingID, _ := m.StoreReceipt(record)

	testCases := map[string]struct {
		inputID            string
		expectedStatusCode int
	}{
		"success": {
			inputID:            existingID,
			expectedStatusCode: http.StatusOK,
		},
		"not found": {
			inputID:            uuid.New().String(),
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for caseName, tc := range testCases {
		t.Run(caseName, func(t *testing.T) {
			res, err := http.Get(srv.URL + fmt.Sprintf(endpointGetPoints, tc.inputID))
			if err != nil {
				t.Errorf("error sending request: %s", err.Error())
				return
			}
			if res.StatusCode != tc.expectedStatusCode {
				t.Errorf("unexpected status code: got %d, want %d", res.StatusCode, tc.expectedStatusCode)
				return
			}
			if tc.expectedStatusCode == http.StatusOK {
				resBody, err := io.ReadAll(res.Body)
				if err != nil {
					t.Errorf("error reading response body: %s", err.Error())
					return
				}
				var pointsResponse entities.PointsResponse
				err = json.Unmarshal(resBody, &pointsResponse)
				if err != nil {
					t.Errorf("error unarmshal response body: %s", err.Error())
					return
				}
				if pointsResponse.Points != record.Points {
					t.Error("incorrect point value in response")
				}
			}
		})
	}
}
