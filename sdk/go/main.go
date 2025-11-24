package main

import (
	"fmt"
	"os"
	"time"

	"storage.client.com/repository"
	"storage.client.com/repository/models"
)


func main(){

	r := repository.NewRepository("http://localhost:4444")
	
	file := "bak.ipa"

	m := &models.Object{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Public: true,
		Deleted: false,
		Filename: "go.sum",
		Extension: ".sum",
		Disk: "C",
		Location: "ya lo averiguare",
		Hash: "No se",


	}

	bytes,err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	uid, err := r.PutObject(*m,bytes,nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("UUID %s",uid)

}