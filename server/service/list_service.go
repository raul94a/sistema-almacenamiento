package service

import (
	"github.com/storage-system/server/models"
	repository "github.com/storage-system/server/repositories"
)

type ListService[T any] interface {
	ListObjects(id T)(*models.Pagination[models.Object], error) 
}

type ListObjectsService[T any] struct {
	Repository *repository.ObjectRepository
}


func (s ListObjectsService[T]) ListObjects(id T, page, limit int) (*models.Pagination[models.Object], error) {
	if limit < 1 || page < 1 {
		return &models.Pagination[models.Object]{}, nil
	}
	return s.Repository.GetUserFiles(id,page,limit)
}