package server

import (
	"net/http"
	"testing"
	"time"
)


func Test_Landing_Page(t *testing.T){
	t.Skip()
	go HttpServer()	
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