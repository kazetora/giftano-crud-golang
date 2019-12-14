package configs

import (
	"fmt"
	"os"
	"log"
	"giftano-crud-golang/database"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func InitDB() *gorm.DB {
	//open a db connection
	var err error
	dbSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	log.Println(fmt.Sprintf("AAAA %s\n", dbSource))
	db, err := gorm.Open(os.Getenv("DB_TYPE"), dbSource)
	db.LogMode(true)
	database.Migrate(db)

	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	// db.AutoMigrate(&todoModel{})

	return db
}