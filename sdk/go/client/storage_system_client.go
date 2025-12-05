package client

import (
	"github.com/google/uuid"
	"storage.client.com/cryptoutils"
	"storage.client.com/repository"
	"storage.client.com/repository/models"
)


type storageSystemClient struct {
	KeyBearer cryptoutils.StorageSystemRsaPublicKey
	Secret *string
	repository *repository.StorageRepository
	Url *string
}

func (c *storageSystemClient) FetchFiles(jwt *string)([]models.Object, error){
	return c.repository.FetchFiles(jwt)
}
func (c *storageSystemClient) GetFile(object models.Object, jwt *string)([]byte, error){
	return c.repository.FetchFile(object,jwt)
}
func (c *storageSystemClient) PutFile(object models.Object, data []byte, jwt *string)(uuid.UUID, error){
	return c.repository.PutObject(object,data,jwt)
}
func (c *storageSystemClient) RemoveFile(object models.Object,jwt *string) error {
	return c.repository.RemoveObject(object,jwt)
}



func NewClient(url string, secret *string) (*storageSystemClient, error){

	secretBearer := &cryptoutils.StorageSystemRsaPublicKey{
		Url: url,
	}
	publicKey, err := secretBearer.GetKey()
	if err != nil {
		return nil,err
	}
	repository := repository.NewRepositoryWithPublicKey(url,publicKey)
	client:= &storageSystemClient{
		KeyBearer: *secretBearer,
		Url: &url,
		repository: repository,
		Secret: secret,
	}
	return client, nil
}