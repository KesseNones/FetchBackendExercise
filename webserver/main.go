package main

import (
	"fmt"
	"net/http"
	"os"
	"errors"
)

func handler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "TEST MESSAGE")
}


func main(){
	http.HandleFunc("/", handler)
	fmt.Println("Listening on port 8000...")
	
	err := http.ListenAndServe(":8000", nil)

	if errors.Is(err, http.ErrServerClosed){
		fmt.Println("Server closed!")
	}else if err != nil{
		fmt.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}
