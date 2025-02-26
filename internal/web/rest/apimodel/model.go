package apimodel

type Product struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Desc      string  `json:"desc"`
	Price     float64 `json:"price"`
	AvgRating float64 `json:"avg_rating"`
}

type Review struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Comment   string `json:"comment"`
	Rating    uint8  `json:"rating"`
	ProductID int    `json:"product_id"`
}
