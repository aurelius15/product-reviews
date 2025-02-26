package apimodel

type Product struct {
	ID        uint32  `json:"id"`
	Name      string  `json:"name"`
	Desc      string  `json:"desc"`
	Price     float64 `json:"price"`
	AvgRating float64 `json:"avg_rating"`
}

type Review struct {
	ID        uint32 `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Comment   string `json:"comment"`
	Rating    uint8  `json:"rating"`
	ProductID uint32 `json:"product_id"`
}
