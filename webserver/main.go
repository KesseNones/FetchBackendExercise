package main

import (
	"fmt"
	"net/http"
	"os"
	"errors"
	"github.com/google/uuid"
)

//Handles root
func Handler(w http.ResponseWriter, r *http.Request){
	fmt.Println(r)
	var newId string = uuid.New().String()
	fmt.Println("TEST", newId)
	fmt.Fprintf(w, "Getting root! Id is: %s", newId)
}

//Handles another endpoint as a test.
func FooHandle(w http.ResponseWriter, r *http.Request){
	fmt.Println(r)
	fmt.Fprintf(w, "Getting foo!")
}

func main(){
	http.HandleFunc("/", Handler)
	http.HandleFunc("/foo", FooHandle) 

	fmt.Println("Listening on Port 8000!")
	
	err := http.ListenAndServe(":8000", nil)

	if errors.Is(err, http.ErrServerClosed){
		fmt.Println("Server closed!")
	}else if err != nil{
		fmt.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}


}
