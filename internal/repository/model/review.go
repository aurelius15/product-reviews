package model

import "time"

type Review struct {
	ID        uint32  `gorm:"primary_key" json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Comment   string  `json:"comment"`
	Rating    uint8   `json:"rating"`
	ProductID uint32  `json:"product_id"`
	Product   Product `json:"product"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (Review) TableName() string {
	return "reviews"
}
