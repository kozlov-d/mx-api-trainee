package common

//Task contains current stats for file processing
type Task struct {
	IsCompleted bool   `json:"isCompleted"`
	TimeSpent   string `json:"timeSpent"` //should be time.Duration, made as string for fancy output
	Created     int    `json:"created"`
	Updated     int    `json:"updated"`
	Missed      int    `json:"missed"`
	Deleted     int    `json:"deleted"`
}

//Merchant represents record in DB
type Merchant struct {
	MerchantID int     `json:"id"`
	Offers     []Offer `json:"offers"`
}

//ContainsOffer returns true if merchant contain offer with given ID
func (m *Merchant) ContainsOffer(oID int) bool {
	for _, n := range m.Offers {
		return n.OfferID == oID
	}
	return false
}

//OfferRow represents a row in .xlsx file.
//Undefined what happens in conversion (int -> float64), if cell actually contains int
type OfferRow struct {
	ID        float64 `xlsx:"0"`
	Name      string  `xlsx:"1"`
	Price     float64 `xlsx:"2"`
	Quantity  float64 `xlsx:"3"`
	Available bool    `xlsx:"4"`
}

//Offer is origin form of row in .xlsx (Google sheets saves number as float). Also can be used as DB model
type Offer struct {
	OfferID    int    `json:"id"`
	OfferName  string `json:"name"`
	Price      int    `json:"price"`
	Quantity   int    `json:"quantity"`
	MerchantID int    `json:"-"`
	Available  bool   `json:"-"`
}

//Convert row to origin offer
func Convert(or OfferRow) *Offer {
	offer := Offer{OfferName: or.Name, Available: or.Available}
	offer.OfferID = int(or.ID)
	offer.Price = int(or.Price)
	offer.Quantity = int(or.Quantity)
	return &offer
}

//Validate fields which passed reading and conversion but might be invalid
func (o Offer) Validate() bool {
	if o.OfferID < 1 {
		return false
	}
	if len(o.OfferName) == 0 {
		return false
	}
	if o.Price < 0 {
		return false
	}
	if o.Quantity < 1 {
		return false
	}

	return true
}
