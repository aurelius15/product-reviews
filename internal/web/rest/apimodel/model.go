package apimodel

type Product struct {
	ID        int     `json:"id"`
	Name      string  `json:"name" validate:"required,min=1,max=255"`
	Desc      string  `json:"desc" validate:"required,min=1,max=1000"`
	Price     float64 `json:"price" validate:"required,gt=0"`
	AvgRating float64 `json:"avg_rating"`
}

type Review struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name" validate:"required,min=1,max=100"`
	LastName  string `json:"last_name" validate:"required,min=1,max=100"`
	Comment   string `json:"comment" validate:"required,min=5,max=500"`
	Rating    uint8  `json:"rating" validate:"required,gte=1,lte=5"`
	ProductID int    `json:"product_id" validate:"required,gte=1"`
}
