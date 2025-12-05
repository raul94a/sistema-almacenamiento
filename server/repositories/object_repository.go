package repository

import (
	"fmt"
	"log"
	"os"

	"github.com/storage-system/server/models"
	"gorm.io/gorm"
)

/**
* Draft for repository
 */

type ObjectRepository struct {
	Db *gorm.DB
}

func (r *ObjectRepository) GetUserFiles(userId any, page, maxItemsPerPage int) (*models.Pagination[models.Object], error) {
	/*
		Discussion: Decision must be taken here.
		A file can be: Owned by an user => where user_id = ?
					   Shared to an user => shared_permissions where user_id = ?
					   Public => I think We should not show these items for sure -but
					   An url can be generated to get them-

		Proposed Query:

		SELECT *
		FROM users u
		INNER JOIN shared_permissions sh
			ON sh.user_id = u.id
		WHERE u.id = $ID
		LIMIT $LIMIT
		OFFSET $OFFSET

		ALSO: Imagine the app is separated into BUCKETS maybe an user only
			  Want to retrieve the files from ONE bucket. We should check also this.

		ALSO: Imagine this is used as a TENANT. This concept has to be engineered
			  more deeply. Here we don't care about the strategy to SAVE a file.
			  We only want to retrieve the objects from a TENANT.

		ALSO: Imagine this is used as a TENANT. This tenant provides some storage
		      to users/groups/organizations/departments in form of BUCKETS

			  Defining a BUCKET: A Storage Space reserved to an user or group of users.
			  					 A BUCKET defines if the Reserved space can be overcome
								 by creating automatically a second bucket or the next
								 uploads MUST be rejected.

			  Defining ownership: An Object is owned by the person who is uploading it,
			                      the available space is substracted from:

								  * Total Storage if user_id is null
								  * Bucket when a group of users share the same bucket -No tenant-
								  	- Total storage
								  * Tenant when an user from a tenant upload a file
								  	- Bucket when the user of a tenant upload to its bucket


	*/
	if maxItemsPerPage < 1 || page < 1 {
		return &models.Pagination[models.Object]{}, nil
	}

	page--
	offset := maxItemsPerPage * page

	var countFiles int64
	if tx := r.Db.Model(&models.Object{}).
		Where("user_owner", userId).
		Count(&countFiles); tx.Error != nil {
		return nil, tx.Error
	}

	var objects []models.Object

	if tx := r.Db.
		Where("user_owner", userId).
		Offset(offset).
		Limit(maxItemsPerPage).
		Find(&objects); tx.Error != nil {
		return nil, tx.Error
	}

	totalPages := (countFiles) / int64(maxItemsPerPage)
	maxItemsPerPage_64 := int64(maxItemsPerPage)

	if countFiles%maxItemsPerPage_64 != 0 {
		totalPages++
	}

	page++
	pagination := &models.Pagination[models.Object]{
		Page:     int(page),
		LastPage: int(totalPages),
		Count:    countFiles,
		Items:    objects,
	}

	return pagination, nil
}

func (r *ObjectRepository) UploadFile(object *models.Object) error {

	if tx := r.Db.Create(&object); tx.Error != nil {
		// TODO: Remove
		log.Printf("UploadFile() %s", tx.Error.Error())
		return tx.Error
	}
	return nil
}

func (r *ObjectRepository) GetUserFile(objectId string) (
	[]byte, error,
) {
	var object models.Object
	if tx := r.Db.
		Where("id = ?", objectId).Find(&object); tx.Error != nil {
		return nil, tx.Error
	}
	// Get User File
	filepath := fmt.Sprintf("%s%s/%s%s", object.Disk, object.Location, object.Filename, object.Extension)
	// large files can be downloaded in chunks?
	//f,e:=os.OpenFile(filepath,1,os.ModeDevice)
	//s, e := f.Stat()
	//s.Size()
	return os.ReadFile(filepath)
}

// Soft delete, gorm automatically handles the deleted_at file
func (r *ObjectRepository) DeleteFile(object *models.Object) error {
	object.Deleted = true
	if tx := r.Db.Delete(object); tx.Error != nil {
		return tx.Error
	}
	return nil
}
