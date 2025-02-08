package entities

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

func (r *Receipt) Validate() bool {
	switch {
	case r.Retailer == "":
		return false
	case r.PurchaseDate == "":
		return false
	case r.PurchaseTime == "":
		return false
	case len(r.Items) == 0:
		return false
	case r.Total == "":
		return false
	}
	return true
}

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

type ReceiptRecord struct {
	Receipt
	Points int
}

type ProcessResponse struct {
	ID string `json:"id"`
}

type PointsResponse struct {
	Points int `json:"points"`
}
