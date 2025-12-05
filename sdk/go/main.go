package main

import (
	"fmt"
	"os"
	"time"

	"storage.client.com/client"
	"storage.client.com/repository"
	"storage.client.com/repository/models"
)


func main(){
	var secret string
	secret = "hola"
	client, err := client.NewClient("HOLA", &secret)
	if client == nil {
		fmt.Sprintf("CLIENT IS NOT INITIALIZED")
		os.Exit(0)
	}
	
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