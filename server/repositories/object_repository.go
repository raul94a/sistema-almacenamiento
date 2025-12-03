package repository

import (
	"context"
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
	ctx context.Context
}


func (r *ObjectRepository) UploadFile(object *models.Object) (error){

	if tx := r.Db.Create(&object); tx.Error != nil {
		// TODO: Remove
		log.Printf("UploadFile() %s",tx.Error.Error())
		return  tx.Error
	}
	return  nil
}

// Soft delete, gorm automatically handles the deleted_at file
func (r *ObjectRepository) DeleteFile(object *models.Object) (error){
	object.Deleted = true
	if tx := r.Db.Delete(object); tx.Error != nil {
		return tx.Error
	}
	return nil;
}