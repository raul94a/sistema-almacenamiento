package models

import (
	"time"
)

type Object struct {
	ID               string     `json:"id"` 
	Description      string     `json:"description"`
	Disk             string     `json:"disk"`
	Location         string     `json:"location"`
	
	Bucket           *string    `json:"bucket"` 
	Region           *string    `json:"region"`
	Parent           *string    `json:"parent"` 

	Encrypted        bool       `json:"encrypted"`
	EncryptionMethod *string    `json:"encryption_method"`
	Public           bool       `json:"public"`
	Deleted          bool       `json:"deleted"` 

	Filename         string     `json:"filename"`
	Extension        string     `json:"extension"`
	Hash             string     `json:"hash"`
	Size             float64    `json:"size"`
	Unit             string     `json:"unit"`

	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at"` 

	UserOwner        string     `json:"user_owner"`
}

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
	UserOwner         *string `json:"user_owner"`
}
