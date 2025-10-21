package service

import "kakebo/internal/repository"

type Service struct {
	DB *repository.Mongo
}

func New(db *repository.Mongo) *Service {
	return &Service{DB: db}
}
