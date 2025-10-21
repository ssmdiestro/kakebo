package dto

type CategoryRecord struct {
	CategoryName  string   `json:"categoryName" bson:"categoryName"`
	CategoryColor string   `json:"categoryColor" bson:"categoryColor"`
	Subcategories []string `json:"subCategories" bson:"subCategories"`
}
