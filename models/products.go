package models

type Products struct {
	Model
	Name string `json:"name"`
	Description string `json:"description"`
}