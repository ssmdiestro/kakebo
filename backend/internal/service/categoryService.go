package service

import (
	"context"
	"kakebo/internal/dto"
	"kakebo/internal/repository"

	"go.mongodb.org/mongo-driver/bson"
)

func (s *Service) GetCategories(ctx context.Context) ([]dto.CategoryRecord, error) {
	coll := s.DB.Collection(repository.GetdatabaseName(), repository.GetcategoryCollection())
	return repository.Find[dto.CategoryRecord](ctx, coll, bson.D{})
}

// func (s *Service) NewCategory(ctx context.Context, categoryRequest dto.CategoryRequest) error {
// 	coll := s.DB.Collection(repository.GetdatabaseName(), repository.GetcategoryCollection())
// 	return repository.InsertOne(ctx, coll, categoryRequest)
// }
