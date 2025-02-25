package model

import "time"

type Product struct {
	ID          uint32   `gorm:"primary_key" json:"id"`
	Name        string   `gorm:"not null" json:"name"`
	Description string   `gorm:"not null" json:"description"`
	Price       float64  `gorm:"not null" json:"price"`
	Reviews     []Review `gorm:"constraint:OnDelete:CASCADE" json:"reviews"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (Product) TableName() string {
	return "products"
}
