package main

import (
	"fmt"
	"net/http"
	"os"
	"errors"
	"github.com/google/uuid"
	"encoding/json"
	"io/ioutil"
)

type DataBase struct {
	names []string
}

type Input struct{
	Name string `json:"name"`
}

//Adds posted name to db.
func (db *DataBase) NameHandler(w http.ResponseWriter, r *http.Request){
	fmt.Println(uuid.New().String())
	if r.Method == http.MethodGet{
		fmt.Println("GET REQUEST")
		if len(db.names) == 0{
			fmt.Fprintf(w, "NO NAMES\n")
		}
		for i := range db.names{
			fmt.Fprintf(w, "%s\n", db.names[i])
		}
	}else if r.Method == http.MethodPost{
		fmt.Println("POST REQUEST")
	
		body, err0 := ioutil.ReadAll(r.Body)
		if err0 != nil{
			fmt.Println("STUFF BROKE HARD")
			return
		}

		var input Input
		
		err := json.Unmarshal(body, &input)
		if err != nil{
			fmt.Println("STUFF BROKE", err)
			return
		}
			
		db.names = append(db.names, input.Name)
	}else{
		fmt.Fprintf(w, "ERROR! Request must be a GET or POST!")
	}
}

func main(){
	data := &DataBase{
		names: []string{},
	}

	http.HandleFunc("/names", data.NameHandler)
	//http.HandleFunc("/names", data.DisplayNames) 

	fmt.Println("Listening on Port 8000!")
	
	err := http.ListenAndServe(":8000", nil)

	if errors.Is(err, http.ErrServerClosed){
		fmt.Println("Server closed!")
	}else if err != nil{
		fmt.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}


}
