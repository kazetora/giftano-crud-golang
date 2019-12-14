package controllers

import (
	"fmt"
	"strings"
	"net/http"
	"encoding/json"
	"strconv"
	"errors"
	"giftano-crud-golang/models"
	"giftano-crud-golang/helpers"
	"giftano-crud-golang/requests"
	"github.com/jinzhu/gorm"
	"github.com/go-chi/chi"
)

type ProductController struct {
	Db *gorm.DB
}

// type ProductResponse struct {
// 	Name string `json:"name"`
// 	Category []string `json:"category"`
// 	Description string `json:"description"`
// }

func NewProductController(db *gorm.DB) *ProductController {
	return &ProductController{
		Db: db,
	}
}

func filterProduct(db *gorm.DB, sql *gorm.DB, r *http.Request) *gorm.DB {
	queries := r.URL.Query()

	for key, _ := range queries {
		switch key {
			case
				"id":
				items := strings.Split(queries.Get(key), ",")
				sql = sql.Where(fmt.Sprintf("%s IN (?)", key),items)
			case
				"keyword":
				sql = sql.Where("name LIKE (\"%" +  queries.Get(key) + "%\") OR description LIKE (\"%" + queries.Get(key) + "%\")")
				
			case
				"category":
				items := strings.Split(queries.Get(key), ",")
				subquery := db.Model(models.ProductCategories{}).Where("category_id IN (?)", items).Select("product_id").SubQuery()
				sql = sql.Where("id IN (?)", subquery)
			case 
				"category_tree":
				categoryRoot := models.Categories{}
				if db.Where("id = ?", queries.Get(key)).First(&categoryRoot).RecordNotFound()  {
					continue
				}
				subquery1 := categoryRoot.GetDecendantsSubQuery(db, "id").SubQuery()
				subquery2 := db.Model(models.ProductCategories{}).Where("category_id IN (?)", subquery1).Select("product_id").SubQuery()
				sql = sql.Where("id IN (?)", subquery2)

		}
	}

	return sql
}

func checkProductID(db *gorm.DB, id int) ([]models.Products, error) {
	var products []models.Products
	sql := db.Model(models.Products{})

	sql.Where("id = ?", id).Find(&products)

	if len(products) != 1 {
		// helpers.RespondError(w, http.StatusBadRequest, "Record not found or found multiple records")
		return products, errors.New("Record not found or found multiple records")
	}

	return products, nil
}

func (p *ProductController) GetProducts(w http.ResponseWriter, r *http.Request) {
	var products []models.Products
	sql := p.Db.Model(models.Products{})
	sql = filterProduct(p.Db, sql, r)
	sql.Find(&products)

	helpers.RespondJSON(w, http.StatusOK, products)
}

func (p *ProductController) RegisterProduct(w http.ResponseWriter, r *http.Request) {
	product := models.Products{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&product); err != nil {
		helpers.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	db := p.Db.Model(models.Products{})
	
	if errs := db.Create(&product).GetErrors(); len(errs) > 0 {
		err := models.CheckErrors(errs, "Create new entry failed")
		helpers.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.RespondJSON(w, http.StatusCreated, product)
}

func (p *ProductController) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	products, err := checkProductID(p.Db, id)
	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&products[0]); err != nil {
		helpers.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	db := p.Db.Model(models.Products{})
	if errs := db.Save(&products[0]).GetErrors(); len(errs) > 0 {
		err := models.CheckErrors(errs, "Update entry failed")
		helpers.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}	

	helpers.RespondJSON(w, http.StatusOK, products[0])
}

func (p *ProductController) UpdateProductCategory(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	_, err := checkProductID(p.Db, id)
	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	param := requests.ProductCategoryID{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&param); err != nil {
		helpers.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	type categoryID struct {
		ID int `json:"id"`
	}
	currentCategoryIds := []categoryID{}
	p.Db.Model(models.ProductCategories{}).Where("product_id = ?", id).Select("category_id as id").Scan(&currentCategoryIds)

	ccid := make([]int, len(currentCategoryIds))

	for _, a := range currentCategoryIds {
		fmt.Printf("%d ", a.ID)
		ccid = append(ccid, a.ID)
	}
	fmt.Println()

	addCategories := helpers.Difference(param.CategoryIDs, ccid)

	delCategories := helpers.Difference(ccid, param.CategoryIDs)

	for _, a := range addCategories {
		fmt.Printf("%d ", a)
	}
	fmt.Println()
	for _, a := range delCategories {
		fmt.Printf("%d ", a)
	}
	fmt.Println()

	// Create
	for _,c := range addCategories {
		item := models.ProductCategories {
			ProductID: id,
			CategoryID: c,
		}
		if errs := p.Db.Create(&item).GetErrors(); len(errs) > 0 {
			err := models.CheckErrors(errs, "Add new category for product failed")
			helpers.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	// Delete
	if len(delCategories) > 0 {
		delCategoryIDs := []categoryID{}
		p.Db.Model(models.ProductCategories{}).Where("product_id = ? AND category_id in (?)", id, delCategories).Select("id").Scan(&delCategoryIDs)
		for _,ct := range delCategoryIDs {
			if errs := p.Db.Where("id = ?", ct.ID).Delete(models.ProductCategories{}).GetErrors() ; len(errs) > 0 {
				err := models.CheckErrors(errs, "Delete old category for product failed")
				helpers.RespondError(w, http.StatusInternalServerError, err.Error())
				return
			}
		}
	}

	helpers.RespondJSON(w, http.StatusOK, map[string][]int{
		"added_categories": addCategories,
		"deleted_categories": delCategories,
	})
}

func (p *ProductController) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	product := models.Products{}
	if errs := p.Db.Where("id = ?", id).First(&product).GetErrors(); len(errs) > 0 {
		err := models.CheckErrors(errs, "Product is not exist")
		helpers.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if errs := p.Db.Delete(&product).GetErrors(); len(errs) > 0 {
		err := models.CheckErrors(errs, "Product delete failed")
		helpers.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if errs := p.Db.Where("product_id = ?", product.ID).Delete(models.ProductCategories{}).GetErrors(); len(errs) > 0 {
		err := models.CheckErrors(errs, "Product category delete failed")
		helpers.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helpers.RespondJSON(w, http.StatusOK, product)
}