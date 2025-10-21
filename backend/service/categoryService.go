package service

import (
	"kakebo/dto"
	"kakebo/repository"
	"os"

	"go.mongodb.org/mongo-driver/bson"
)

func GetCategories() []dto.CategoryRecord {
	categories, err := repository.SelectWithError([]dto.CategoryRecord{}, bson.M{}, os.Getenv("CATEGORY_TABLE"), os.Getenv("KAKEBO_DB"))
	if err != nil {
		return []dto.CategoryRecord{}
	}
	return categories
}
