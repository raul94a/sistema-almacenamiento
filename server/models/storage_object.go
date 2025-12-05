package models

import (
    "time"

    "gorm.io/gorm"
)

type Object struct {
    // Primary Key
    ID string `json:"id" gorm:"primaryKey;size:36"` // assuming UUID

    Description string `json:"description" gorm:"column:description"`
    Disk        string `json:"disk" gorm:"column:disk"`
    Location    string `json:"location" gorm:"column:location"`

    // Optional string fields (NULLable in DB)
    Bucket           *string `json:"bucket,omitempty" gorm:"column:bucket"`
    Region           *string `json:"region,omitempty" gorm:"column:region"`
    Parent           *string `json:"parent,omitempty" gorm:"column:parent;index"` // index for tree queries
    EncryptionMethod *string `json:"encryption_method,omitempty" gorm:"column:encryption_method"`

    // Booleans
    Encrypted bool `json:"encrypted" gorm:"default:false;column:encrypted"`
    Public    bool `json:"public" gorm:"default:false;column:public;index"` // index if filtering public objects
    Deleted   bool `json:"deleted" gorm:"default:false;column:deleted;index"`

    // File metadata
    Filename  string `json:"filename" gorm:"column:filename;not null"`
    Extension string `json:"extension" gorm:"column:extension"`
    Hash      string `json:"hash" gorm:"column:hash;uniqueIndex:idx_hash_deleted"` // prevent duplicates
    Size      int64  `json:"size" gorm:"column:size"` // ALWAYS use int64 for bytes!
    Unit      string `json:"unit" gorm:"column:unit"` // e.g., "bytes", "KB" – consider removing if Size is always in bytes

    // Timestamps
    CreatedAt time.Time      `json:"created_at" gorm:"column:created_at"`
    UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at"`
    DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"` // Proper soft delete

    // Foreign key
    UserOwner *string `json:"user_owner" gorm:"column:user_owner;index"`
    // If you have a User model:
    // UserOwner   string `json:"-" gorm:"column:user_owner;size:36;not null;index"`
    // User        User   `json:"user,omitempty" gorm:"foreignKey:UserOwner;references:ID"`
}