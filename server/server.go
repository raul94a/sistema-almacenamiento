package server

import (
	"net/http"
	"os"
)

var indexHtmlCache []byte		

func HttpServer(){

	publicDir := "./server/public"

	fs := http.FileServer(http.Dir(publicDir))

    http.Handle("/public/",http.StripPrefix("/public/",fs))

	http.HandleFunc("/index.html", func(w http.ResponseWriter, r *http.Request) {
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

	http.ListenAndServe(":4444",nil)
}