package repository

import (
	"github.com/google/uuid"
	"storage.client.com/repository/models"
)

type Repository interface {
	FetchFiles(jwt *string) ([]models.Object, error)

	FetchFile(o models.Object, jwt *string) ([]byte, error)

	PutObject(object models.Object, jwt *string) (uuid.UUID, error)

	RemoveObject(o models.Object, jwt *string) error

}