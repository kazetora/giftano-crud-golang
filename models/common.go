package models

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
)
func CheckErrors(errs []error, message string) error {
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Println(err)
		}
		return errors.New(message)
	}

	return nil
}

func AddItem(db *gorm.DB, data interface{}) error{
	errs := db.Create(&data).GetErrors()

	return CheckErrors(errs, "Create new entry failed")
}

func UpdateItem(db *gorm.DB, id int, data interface{}) error {
	var count int
	db.Where("id = ?", id).Count(&count)
	if count < 1 {
		return errors.New("Record is not found")
	}

	errs := db.Save(&data).GetErrors()

	return CheckErrors(errs,"Update entry failed" )
}

func DeleteItem(db *gorm.DB, id int) error{
	// var count int
	var record []interface{}
	db.Where("id = ?", id).Find(&record)
	if len(record) < 1 {
		return errors.New("Record is not found")
	}
	errs := db.Delete(&record[0]).GetErrors()
	
	return CheckErrors(errs, "Delete entry failed")
}