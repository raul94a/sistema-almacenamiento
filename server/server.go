package server

import (
	"net/http"
	"os"

	"github.com/storage-system/server/controller"
	"github.com/storage-system/server/power"
)

var indexHtmlCache []byte

func HttpServer(uoc controller.UploadObjectController){

	publicDir := "./server/public"

	fs := http.FileServer(http.Dir(publicDir))

    http.Handle("/public/",http.StripPrefix("/public/",fs))


	power.Handler("/index.html",func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type","text/html")
		if indexHtmlCache != nil {

			 w.Write(indexHtmlCache)
			
			 return
		}
		bytes, err := os.ReadFile("./server/public/index.html")
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		indexHtmlCache = bytes
		w.Write(bytes)
		fs.ServeHTTP(w, r)

	})


	power.Handler("/api/v1/PutObject", uoc.UploadObject)


	http.ListenAndServe(":4444",nil)
}