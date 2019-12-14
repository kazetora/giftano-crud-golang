package controllers

import (
	// "fmt"
	// "strings"
	"net/http"
	"encoding/json"
	"strconv"
	"giftano-crud-golang/models"
	"giftano-crud-golang/helpers"
	"github.com/jinzhu/gorm"
	"github.com/go-chi/chi"
)

type CategoryController struct {
	Db *gorm.DB
}

func NewCategoryController(db *gorm.DB) *CategoryController {
	return &CategoryController{
		Db: db,
	}
}

type CategoryTreeNode struct {
	Name string `json:"name"`
	ID int `json:"id"`
	Children []*CategoryTreeNode `json:"children;omitempty"`
}

const (
	moveOnAddSql = "UPDATE categories SET categories.left=CASE WHEN categories.left > ? THEN categories.left+2 ELSE categories.left END, categories.right=CASE WHEN categories.right > ? THEN categories.right+2 ELSE categories.right END"
	moveOnDeleteSql = "UPDATE categories SET categories.left=CASE WHEN categories.left > ? THEN categories.left-? ELSE categories.left END, categories.right=CASE WHEN categories.right > ? THEN categories.right-? ELSE categories.right END"
	moveOnLevelUpSql = "UPDATE categories SET categories.left = categories.left-1, categories.right = right-1, categories.depth = categories.depth-1 WHERE categories.left between ? AND ?"
	updateParentIDSql = "UPDATE categories AS child, categories AS parent SET child.parent_id = parent.parent_id WHERE child.parent_id = parent.id AND child.left BETWEEN ? AND ?"
)

func (c *CategoryController) AddCategoryRoot(w http.ResponseWriter, r *http.Request) {
	c.Db.Exec(moveOnAddSql, 0, 0)

	newCategory := models.Categories{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newCategory); err != nil {
		helpers.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()	

	//  0, 1, 1, 2
	newCategory.ParentID = 0
	newCategory.Depth = 1
	newCategory.Left = 1
	newCategory.Right = 2

	db := c.Db.Model(models.Categories{})

	if errs := db.Create(&newCategory).GetErrors(); len(errs) > 0 {
		err := models.CheckErrors(errs, "Create new entry failed")
		helpers.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.RespondJSON(w, http.StatusCreated, newCategory)
}

func (c *CategoryController) AddCategoryByParent(w http.ResponseWriter, r *http.Request) {
	parentID, _ := strconv.Atoi(chi.URLParam(r, "id"))

	parentCategory := models.Categories{}
	newCategory := models.Categories{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newCategory); err != nil {
		helpers.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()	

	db := c.Db.Model(models.Categories{}).Where(" id = ? ", parentID)
	
	if errs := db.First(&parentCategory).GetErrors(); len(errs) > 0 {
		err := models.CheckErrors(errs, "Create new entry failed")
		helpers.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}	

	c.Db.Exec(moveOnAddSql, parentCategory.Right, parentCategory.Right-1)


	//  parentID, parentDepth + 1, parentRight, parentRight + 1
	newCategory.ParentID = parentID
	newCategory.Depth = parentCategory.Depth + 1
	newCategory.Left = parentCategory.Right 
	newCategory.Right = parentCategory.Right + 1

	db = c.Db.Model(models.Categories{})

	if errs := db.Create(&newCategory).GetErrors(); len(errs) > 0 {
		err := models.CheckErrors(errs, "Create new entry failed")
		helpers.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.RespondJSON(w, http.StatusCreated, newCategory)

}

func (c *CategoryController) AddCategoryBySibling(w http.ResponseWriter, r *http.Request) {
	siblingID, _ := strconv.Atoi(chi.URLParam(r, "id"))

	siblingCategory := models.Categories{}
	newCategory := models.Categories{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newCategory); err != nil {
		helpers.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	
	db := c.Db.Model(models.Categories{}).Where(" id = ? ", siblingID)
	
	if errs := db.First(&siblingCategory).GetErrors(); len(errs) > 0 {
		err := models.CheckErrors(errs, "Create new entry failed")
		helpers.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	c.Db.Exec(moveOnAddSql, siblingCategory.Right, siblingCategory.Right)


	//  parentID, siblingDepth, siblingRight + 1, siblingRight + 2
	newCategory.ParentID = siblingCategory.ParentID
	newCategory.Depth = siblingCategory.Depth
	newCategory.Left = siblingCategory.Right  + 1
	newCategory.Right = siblingCategory.Right + 2

	db = c.Db.Model(models.Categories{})

	if errs := db.Create(&newCategory).GetErrors(); len(errs) > 0 {
		err := models.CheckErrors(errs, "Create new entry failed")
		helpers.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.RespondJSON(w, http.StatusCreated, newCategory)
}

func insertChildren(db *gorm.DB, root *CategoryTreeNode) {
	children := []models.Categories{}
	db.Where("parent_id = ?", root.ID).Find(&children)
	if len(children) < 1 {
		return
	}
	root.Children = make([]*CategoryTreeNode, 0, len(children))
	for _,child  := range children {
		newChild := CategoryTreeNode {
			Name: child.CategoryName,
			ID: int(child.ID),
		}
		root.Children = append(root.Children, &newChild)
		insertChildren(db, &newChild)
	}
}

func (c *CategoryController) GetAllCategoryTree(w http.ResponseWriter, r *http.Request) {
	root := CategoryTreeNode{
		Name: "root",
		ID: 0,
	}

	dataModel := models.Categories{}

	firstChildren, err := dataModel.GetAllRoot(c.Db)
	if err != nil {
		helpers.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	root.Children = make([]*CategoryTreeNode, 0, len(firstChildren))
	for _,child := range firstChildren {
		newChild := CategoryTreeNode{
			Name: child.CategoryName,
			ID: int(child.ID),
		}
		root.Children = append(root.Children, &newChild)
		insertChildren(c.Db, &newChild)
	}

	helpers.RespondJSON(w, http.StatusCreated, root.Children)
}

func (c *CategoryController) GetCategoryTreeFromId(w http.ResponseWriter, r *http.Request) {
	rootID, _ := strconv.Atoi(chi.URLParam(r, "id"))
	targetCategory := models.Categories{}

	firstChildren, err := targetCategory.GetAllChildren(c.Db, rootID)
	if err != nil {
		helpers.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	c.Db.Where("id = ?", rootID).First(&targetCategory)

	root := CategoryTreeNode{
		ID: int(targetCategory.ID),
		Name: targetCategory.CategoryName,
	}

	root.Children = make([]*CategoryTreeNode, 0, len(firstChildren))
	for _,child := range firstChildren {
		newChild := CategoryTreeNode{
			Name: child.CategoryName,
			ID: int(child.ID),
		}
		root.Children = append(root.Children, &newChild)
		insertChildren(c.Db, &newChild)
	}

	helpers.RespondJSON(w, http.StatusCreated, root)
}

func (c *CategoryController) RemoveCategorySubtree(w http.ResponseWriter, r *http.Request) {
	rootID, _ := strconv.Atoi(chi.URLParam(r, "id"))
	targetCategory := models.Categories{}

	if errs := c.Db.Where("id = ?", rootID).First(&targetCategory).GetErrors(); len(errs) > 0 {
		err := models.CheckErrors(errs, "Can not find category for delete")
		helpers.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	width := targetCategory.Right - targetCategory.Left

	tx := c.Db.Begin()

	if errs := tx.Where("categories.left BETWEEN ? AND ?", targetCategory.Left, targetCategory.Right).Delete(&models.Categories{}).GetErrors(); len(errs) > 0 {
		tx.Rollback()
		err := models.CheckErrors(errs, "Can not find category for delete")
		helpers.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if errs := tx.Exec(moveOnDeleteSql, targetCategory.Left, width, targetCategory.Right, width).GetErrors(); len(errs) > 0 {
		tx.Rollback()
		err := models.CheckErrors(errs, "Can not update ")
		helpers.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	tx.Commit()

	helpers.RespondJSON(w, http.StatusCreated, targetCategory)
}

func (c *CategoryController) RemoveOneCategory(w http.ResponseWriter, r *http.Request) {
	rootID, _ := strconv.Atoi(chi.URLParam(r, "id"))
	targetCategory := models.Categories{}

	if errs := c.Db.Where("id = ?", rootID).First(&targetCategory).GetErrors(); len(errs) > 0 {
		err := models.CheckErrors(errs, "Can not find category for delete")
		helpers.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	tx := c.Db.Begin()

	if errs := tx.Exec(updateParentIDSql, targetCategory.Left, targetCategory.Right).GetErrors(); len(errs) > 0 {
		tx.Rollback()
		err := models.CheckErrors(errs, "Can not update parent id")
		helpers.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	if errs := tx.Where("id = ?", targetCategory.ID).Delete(&models.Categories{}).GetErrors(); len(errs) > 0 {
		tx.Rollback()
		err := models.CheckErrors(errs, "Can not delete category")
		helpers.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if errs := tx.Exec(moveOnLevelUpSql, targetCategory.Left, targetCategory.Right).GetErrors(); len(errs) > 0 {
		tx.Rollback()
		err := models.CheckErrors(errs, "Can not update parent id")
		helpers.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if errs := tx.Exec(moveOnDeleteSql, targetCategory.Right, 2, targetCategory.Right, 2).GetErrors(); len(errs) > 0 {
		tx.Rollback()
		err := models.CheckErrors(errs, "Can not update parent id")
		helpers.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	tx.Commit()

	helpers.RespondJSON(w, http.StatusCreated, targetCategory)
}