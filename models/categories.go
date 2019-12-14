package models

import (
	"fmt"
	"strings"
	"github.com/jinzhu/gorm"
)

type Categories struct {
	Model
	CategoryName string `gorm:"not null" valid:"required" json:"category_name"`
	ParentID int `gorm:"not null" valid:"required" json:"parent_id"`
	Depth int `gorm:"not null;index:depth" valid:"required" json:"depth"`
	Left int `gorm:"not null;index:leftidx" valid:"required" json:"left"`
	Right int `gorm:"not null;index:rightidx" valid:"required" json:"right"`
}

func (c *Categories) GetAllRoot(db *gorm.DB) ([]Categories, error) {
	return c.GetAllChildren(db, 0)
}

func (c *Categories) GetAllChildren(db *gorm.DB, id int) ([]Categories, error) {
	res := []Categories{}
	if err := CheckErrors(db.Where("parent_id = ?", id).Order("categories.left ASC").Find(&res).GetErrors(), "Something is not right"); err != nil {
		return res, err
	}
	return res, nil
}

func (c *Categories) GetAllRooSubQuery(db *gorm.DB, selects ...string) *gorm.DB {

	return db.Model(c).Where("parent_id = 0").Select(strings.Join(selects, ", ")).Order("categories.left ASC")
}

func (c *Categories) GetDecendantsSubQuery(db *gorm.DB, selects ...string) *gorm.DB {
	s := make([]string, 0, len(selects))
	for _,item := range selects {
		s = append(s, fmt.Sprintf("child.%s", item))
	}
	return db.Raw(fmt.Sprintf("SELECT %s FROM categories AS child, categories AS parent WHERE parent.id = %d AND child.left BETWEEN parent.left AND parent.right ORDER BY child.left ASC", strings.Join(s,", "), c.ID))
}

func (c *Categories) GetNodesByDepth(db *gorm.DB, depth int, selects ...string) *gorm.DB {
	return db.Model(c).Where("depth = ?", depth).Select(strings.Join(selects, ", "))
}