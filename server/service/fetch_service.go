package service

import repository "github.com/storage-system/server/repositories"

type FetchService interface {
	GetObject(objectId string)([]byte,error)
}


type GetObjectService struct {
	Repository *repository.ObjectRepository
}


func (s GetObjectService) GetObject(objectId string) (
	[]byte, error,
) {
	return s.Repository.GetUserFile(objectId)
}