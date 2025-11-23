package models

import "time"



type StorageUnit string

const (
	Bit  StorageUnit  = "bits"
	Byte StorageUnit  = "byte"
	Kb	 StorageUnit  = "Kb"
	Mb	 StorageUnit  = "Mb"
	Gb   StorageUnit  = "Gb"
)



type StorageObjectMetadata struct {
	Filename string
	Size float64
	Unit StorageUnit
	CreatedAt time.Duration
	ModifiedAt time.Duration
	Extension string
}

// this is the main object
type StorageObject struct {
	Id string
	Bucket *string	// Definitely the Bucket is a virtual space inside a Filesystem
					// To Write/Read into a bucket, the user that is sending / accessing the file
					// needs to have Write and/or Read permissions inside the bucket.
					// This idea is useful when a we want to have two different apps
					// that cannot access the same virtual space.
					// A bucket is attached to an App. So maybe
					// will be necesary to explore the Idea of AppId as a
					// way of communicating with a bucket.
	Description string 
	Disk string     // the disk to be stored
	Encrypted bool
	EncryptionMethod *string // Reponsability of the client. The idea is to have 
							// a E2EE system.
	Location string  // where the file will be stored inside the filesystem
	Region *string  // for distributed systems we need to know the region. ej: aws-west-2
	Parent *string // A parent is a directory where the file is stored
	Public bool   // Can any user access this file? true if public access false if private
	Data []byte  // Nothing to comment.
	Metadata StorageObjectMetadata
	
}