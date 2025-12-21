package server

import (
	"net/http"
	"testing"
	"time"

	"github.com/storage-system/server/controller"
	repository "github.com/storage-system/server/repositories"
	"github.com/storage-system/server/service"
)


func Test_Landing_Page(t *testing.T){
	t.Skip()
	uoc := controller.UploadObjectController {
		UploadService: &service.UploadObjectService{
			Repository: &repository.ObjectRepository{},
		},
	}
	go HttpServer(uoc)	
	time.Sleep(5 *time.Second)
	if r,e := http.Get("http://localhost:4444/index.html"); e != nil {
		t.Fatal(e.Error())
	} else {
		var bytes []byte
		_,e = r.Body.Read(bytes)
		if e!=nil{
			t.Fatal(e.Error())
		}
		if r.StatusCode != 200 {
			t.Fatalf("StatusCode %d for body %s",r.StatusCode,string(bytes))
		}
	}


}