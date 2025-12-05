package models

import "time"

type UploadObject struct {
	Description 	  string  `json:"description"`
	Bucket 			  *string `json:"bucket,omitempty"`
	Parent            *string `json:"parent,omitempty"`
    EncryptionMethod  *string `json:"encryption_method,omitempty"`
	Public    		  bool    `json:"public"`
    Filename  		  string  `json:"filename"`
    Extension 	      string  `json:"extension"`
    Hash      	      string  `json:"hash"` 
    Size      		  int64   `json:"size"`
    Unit      		  string  `json:"unit"`
	UserOwner         *string `json:"user_owner" gorm:"column:user_owner;index"`
}


func (o *UploadObject) MapToObject(id string) Object {
	now := time.Now()

	return Object{
		ID: 				id,
        Description:      	o.Description,
        Bucket:           	o.Bucket,
        Parent:           	o.Parent,
        EncryptionMethod: 	o.EncryptionMethod,

        Encrypted: 			o.EncryptionMethod != nil,
        Public:    			o.Public,
        Deleted:   			false,

        Filename:  			o.Filename,
        Extension: 			o.Extension,
        Hash:      			o.Hash,
        Size:      			o.Size,
        Unit:      			"bytes", 

        CreatedAt: 			now,
        UpdatedAt: 			now,
		// TODO: Comprobation of ownership
        UserOwner: 			o.UserOwner, 
    }
}