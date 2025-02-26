package model

import "time"

type Review struct {
	ID        int     `gorm:"primaryKey"`
	FirstName string  `gorm:"not null"`
	LastName  string  `gorm:"not null"`
	Comment   string  `gorm:"not null"`
	Rating    uint8   `gorm:"not null"`
	ProductID int     `gorm:"not null"`
	Product   Product `gorm:"foreignKey:ProductID;references:ID;constraint:OnDelete:CASCADE"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (Review) TableName() string {
	return "reviews"
}
