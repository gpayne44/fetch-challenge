package process

import (
	"math"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gpayne44/fetch-challenge/internal/entities"
)

const (
	pointValueEvenDollar      = 50
	pointValueQuarterMultiple = 25
	pointValueTwoItems        = 5
	pointValueOddPurchaseDate = 6
	pointValueHappyHours      = 10

	dateFmt = "2006-01-02"
	timeFmt = "15:04"
)

func CalculatePoints(receipt entities.Receipt) (int, []error) {
	var (
		totalPoints int
		errors      []error
	)
	totalPoints += calculateNamePoints(receipt.Retailer)

	pricePoints, err := calculateTotalPricePoints(receipt.Total)
	if err != nil {
		errors = append(errors, err)
	}
	totalPoints += pricePoints

	totalPoints += calculateItemsPoints(len(receipt.Items))

	descriptionPoints, err := calculateDescriptionPoints(receipt.Items)
	if err != nil {
		errors = append(errors, err)
	}
	totalPoints += descriptionPoints

	oddDatePoints, err := calculateOddDatePoints(receipt.PurchaseDate)
	if err != nil {
		errors = append(errors, err)
	}
	totalPoints += oddDatePoints

	happyHoursPoints, err := calculateHappyHoursPoints(receipt.PurchaseTime)
	if err != nil {
		errors = append(errors, err)
	}
	totalPoints += happyHoursPoints

	return totalPoints, errors
}

// one point for every alphanumeric character in the retailer name
func calculateNamePoints(retailerName string) int {
	if retailerName == "" {
		return 0
	}
	var namePoints int
	nameBytes := []byte(retailerName)
	for _, b := range nameBytes {
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') {
			namePoints++
		}
	}
	return namePoints
}

// 50 points if the total is a round dollar amount with no cents
// 25 points if the total is a multiple of 0.25
func calculateTotalPricePoints(total string) (int, error) {
	receiptTotal, err := strconv.ParseFloat(total, 64)
	if err != nil {
		return 0, err
	}
	var totalPricePoints int
	if receiptTotal == math.Trunc(receiptTotal) {
		totalPricePoints += pointValueEvenDollar
	}

	if rem := math.Mod(receiptTotal, 0.25); rem == 0 {
		totalPricePoints += pointValueQuarterMultiple
	}

	return totalPricePoints, nil
}

// 5 points for every two items on the receipt.
func calculateItemsPoints(itemLen int) int {
	numberOfItemPairs := itemLen / 2
	return numberOfItemPairs * pointValueTwoItems
}

// If the trimmed length of the item description is a multiple of 3,
// multiply the price by 0.2 and round up to the nearest integer.
// The result is the number of points earned.
func calculateDescriptionPoints(items []entities.Item) (int, error) {
	var descriptionPoints int
	if len(items) == 0 {
		return 0, nil
	}
	for _, item := range items {
		trimmed := strings.TrimSpace(item.ShortDescription)
		charCount := utf8.RuneCountInString(trimmed)
		if charCount%3 == 0 {
			priceFloat, err := strconv.ParseFloat(item.Price, 64)
			if err != nil {
				return 0, err
			}
			descriptionPoints += int(math.Ceil(priceFloat * 0.2))
		}
	}

	return descriptionPoints, nil
}

// 6 points if the day in the purchase date is odd.
func calculateOddDatePoints(purchaseDate string) (int, error) {
	parsedDate, err := time.Parse(dateFmt, purchaseDate)
	if err != nil {
		return 0, err
	}
	dayVal := parsedDate.Day()
	if dayVal%2 != 0 {
		return pointValueOddPurchaseDate, nil
	}
	return 0, nil
}

// 10 points if the time of purchase is after 2:00pm and before 4:00pm.
func calculateHappyHoursPoints(purchaseTime string) (int, error) {
	parsedTime, err := time.Parse(timeFmt, purchaseTime)
	if err != nil {
		return 0, err
	}
	hour := parsedTime.Hour()
	if hour == 14 || hour == 15 {
		return pointValueHappyHours, nil
	}
	return 0, nil
}
