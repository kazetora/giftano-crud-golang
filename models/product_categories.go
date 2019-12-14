package models

type ProductCategories struct {
	Model
	ProductID int `gorm:"not null" valid:"required" json:"product_id"`
	CategoryID int `gorm:"not null" valid:"required" json:"category_id"`
}