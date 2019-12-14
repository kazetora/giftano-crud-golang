package routers

import (
	"giftano-crud-golang/configs"
	"giftano-crud-golang/controllers"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jinzhu/gorm"
)

func APIRouter(db *gorm.DB) *chi.Mux {
	r := chi.NewRouter()

	ProductController := controllers.NewProductController(db)
	CategoryController := controllers.NewCategoryController(db)

	r.Route("/product", func(r chi.Router) {
		r.Get("/", ProductController.GetProducts)
		r.Post("/", ProductController.RegisterProduct)
		r.Put("/{id:[0-9]+}", ProductController.UpdateProduct)
		r.Post("/updateCategory/{id:[0-9]+}", ProductController.UpdateProductCategory)
		r.Delete("/{id:[0-9]+}", ProductController.DeleteProduct)
	})

	r.Route("/category", func(r chi.Router) {
		r.Get("/getCategoryTree", CategoryController.GetAllCategoryTree)
		r.Get("/getCategoryTreeFromId/{id:[0-9]+}", CategoryController.GetCategoryTreeFromId)
		r.Post("/addRoot", CategoryController.AddCategoryRoot)
		r.Post("/addByParentId/{id:[0-9]+}", CategoryController.AddCategoryByParent)
		r.Post("/addBySiblingId/{id:[0-9]+}", CategoryController.AddCategoryBySibling)
		r.Delete("/deleteCategorySubtree/{id:[0-9]+}", CategoryController.RemoveCategorySubtree)
		r.Delete("/RemoveOneCategory/{id:[0-9]+}", CategoryController.RemoveCategorySubtree)
	})

	return r
}

func InitRouter() *chi.Mux {
	r := chi.NewRouter()

	db := configs.InitDB()

	
	
	r.Use(
		middleware.Recoverer,
		middleware.Logger,
	)

	r.Mount("/api/v1", APIRouter(db))

	return r
}