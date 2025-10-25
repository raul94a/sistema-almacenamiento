package api


// TODO: things

type FileRepository interface {
	ListObjects()
	GetObject()
	PutObject()
	DeleteObject()
}