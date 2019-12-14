package database

import (
	"github.com/jinzhu/gorm"
	"giftano-crud-golang/models"
)

func Migrate(db *gorm.DB) {
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&models.Products{}, 
										&models.Categories{},
										&models.ProductCategories{})
}