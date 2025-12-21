package main

import (
	"fmt"
	"os"

	"storage.client.com/client"
	"storage.client.com/repository"
	"storage.client.com/repository/models"
)


func main(){
	fmt.Println("INIT")
	var secret string
	secret = "hola"
	client, err := client.NewClient("HOLA", &secret)
	if client == nil {
		fmt.Sprintf("CLIENT IS NOT INITIALIZED")
		os.Exit(0)
	}
	
	r := repository.NewRepository("http://localhost:4444")
	
	file := "DATA.txt"

	fmt.Println("SET FILE")

	m := &models.UploadObject{
		Description: "Desc",	 
		Parent         : nil,
		EncryptionMethod: nil,
		Public    		:true,
		Filename  		: "DATA.txt",
		Extension 	     : ".txt",
		Hash      	      :"HASH",
		Size      		  : 22,
		Unit      		  :"Bytes",
		UserOwner         : nil,
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