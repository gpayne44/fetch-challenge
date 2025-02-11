package process

import (
	"testing"

	"github.com/gpayne44/fetch-challenge/internal/entities"
)

func Test_calculateNamePoints(t *testing.T) {
	testCases := map[string]struct {
		inputName     string
		expectedScore int
	}{
		"basic success": {
			inputName:     "Target",
			expectedScore: 6,
		},
		"success with spaces": {
			inputName:     "Fred Meyer",
			expectedScore: 9,
		},
		"non-alphanumeric counters are not counted": {
			inputName:     "T & T Supermarket",
			expectedScore: 13,
		},
	}

	for caseName, tc := range testCases {
		t.Run(caseName, func(t *testing.T) {
			nameScore := calculateNamePoints(tc.inputName)
			if nameScore != tc.expectedScore {
				t.Errorf("nameScore is not correct: got %d, want %d", tc.expectedScore, nameScore)
			}
		})
	}
}

func Test_calculateTotalPricePoints(t *testing.T) {
	testCases := map[string]struct {
		inputPrice    string
		expectedScore int
		expectError   bool
	}{
		"quarter multiple": {
			inputPrice:    "35.25",
			expectedScore: 25,
		},
		"whole dollar amount is always quarter multiple": {
			inputPrice:    "10.00",
			expectedScore: 75,
		},
		"price is not a valid number": {
			inputPrice:  "asdf.00",
			expectError: true,
		},
	}

	for caseName, tc := range testCases {
		t.Run(caseName, func(t *testing.T) {
			priceScore, err := calculateTotalPricePoints(tc.inputPrice)
			if tc.expectError {
				if err == nil {
					t.Error("expected error but did not get one")
				}
			}
			if priceScore != tc.expectedScore {
				t.Errorf("nameScore is not correct: got %d, want %d", priceScore, tc.expectedScore)
			}
		})
	}
}

func Test_calculateItemsPoints(t *testing.T) {
	testCases := []struct {
		inputItemsCount int
		expectedScore   int
	}{
		{
			inputItemsCount: 10,
			expectedScore:   25,
		},
		{
			inputItemsCount: 15,
			expectedScore:   35,
		},
		{
			inputItemsCount: 7,
			expectedScore:   15,
		},
		{
			inputItemsCount: 23,
			expectedScore:   55,
		},
	}

	for _, tc := range testCases {
		score := calculateItemsPoints(tc.inputItemsCount)
		if score != tc.expectedScore {
			t.Errorf("unexpected score: got %d, want %d", score, tc.expectedScore)
		}
	}
}

func Test_calculateDescriptionPoints(t *testing.T) {
	testCases := map[string]struct {
		inputItems    []entities.Item
		expectedScore int
		expectError   bool
	}{
		"2 scoring descriptions": {
			inputItems: []entities.Item{
				{
					ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ",
					Price:            "12.00",
				},
				{
					ShortDescription: "Emils Cheese Pizza",
					Price:            "12.25",
				},
				{
					ShortDescription: "Mountain Dew 12PK",
					Price:            "6.49",
				},
			},
			expectedScore: 6,
		},
		"nothing scoring": {
			inputItems: []entities.Item{
				{
					ShortDescription: "Gatorade",
					Price:            "2.25",
				},
				{
					ShortDescription: "Gatorade",
					Price:            "2.25",
				},
			},
			expectedScore: 0,
		},
		"invalid prices": {
			inputItems: []entities.Item{
				{
					ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ",
					Price:            "12.asdf",
				},
				{
					ShortDescription: "Emils Cheese Pizza",
					Price:            ";lkj.25",
				},
				{
					ShortDescription: "Mountain Dew 12PK",
					Price:            "6.49",
				},
			},
			expectedScore: 0,
			expectError:   true,
		},
	}

	for _, tc := range testCases {
		score, err := calculateDescriptionPoints(tc.inputItems)
		if tc.expectError && err == nil {
			t.Errorf("expected error but did not get one")
		}
		if score != tc.expectedScore {
			t.Errorf("unexpected score: got %d, want %d", score, tc.expectedScore)
		}
	}
}

func Test_calculateOddDatePoints(t *testing.T) {
	testCases := map[string]struct {
		inputDate     string
		expectedScore int
		expectError   bool
	}{
		"valid odd date": {
			inputDate:     "2022-01-01",
			expectedScore: pointValueOddPurchaseDate,
		},
		"valid even date": {
			inputDate: "2022-03-20",
		},
		"invalid date": {
			inputDate:   "202234-123-234-908",
			expectError: true,
		},
	}

	for caseName, tc := range testCases {
		t.Run(caseName, func(t *testing.T) {
			score, err := calculateOddDatePoints(tc.inputDate)
			if tc.expectError && err == nil {
				t.Errorf("expected error but did not get one")
			}
			if score != tc.expectedScore {
				t.Errorf("unexpected score: got %d, want %d", score, tc.expectedScore)
			}
		})
	}
}

func Test_calculateHappyHoursPoints(t *testing.T) {
	testCases := map[string]struct {
		inputTime     string
		expectedScore int
		expectError   bool
	}{
		"scoring time": {
			inputTime:     "14:33",
			expectedScore: pointValueHappyHours,
		},
		"non-scoring time": {
			inputTime: "13:01",
		},
		"invalid time": {
			inputTime:   "01:",
			expectError: true,
		},
	}

	for caseName, tc := range testCases {
		t.Run(caseName, func(t *testing.T) {
			score, err := calculateHappyHoursPoints(tc.inputTime)
			if tc.expectError && err == nil {
				t.Errorf("expected error but did not get one")
			}
			if score != tc.expectedScore {
				t.Errorf("unexpected score: got %d, want %d", score, tc.expectedScore)
			}
		})
	}
}
