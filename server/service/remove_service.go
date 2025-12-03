package service

import 
(

	"github.com/storage-system/server/models"
	repo "github.com/storage-system/server/repositories"

)

type RemoveService interface {
	Delete(object models.Object) (error)
	
}

type RemoveObjectService struct {
	Repository *repo.ObjectRepository
}

func (s *RemoveObjectService) Delete(object *models.Object) error {
	return s.Repository.DeleteFile(object)
}