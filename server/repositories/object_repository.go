package repository

import (
	"log"
	"github.com/storage-system/server/models"
	"gorm.io/gorm"
)

/**
* Draft for repository
*
*
*
*
*/


type ObjectRepository struct {
	Db *gorm.DB
}

func (r *ObjectRepository) UploadFile(object *models.Object) (error){

	if tx := r.Db.Create(&object); tx.Error != nil {
		// TODO: Remove
		log.Printf("UploadFile() %s",tx.Error.Error())
		return  tx.Error
	}
	return  nil
}