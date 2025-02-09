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

type ReceiptItem struct{
	ShortDesc string `json:"shortDescription"`
	Price string `json:"price"`
}

//Input struct used to handle either receipt input or just id.
type Receipt struct{
	Retailer string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	Total string `json:"total"`
	Items []ReceiptItem `json:"items"`
}

type DataBase struct {
	Data map[string]int
}

//Takes in the input receipt, calculates point value, 
// and writes points to database.
func (db *DataBase) InsertToDatabase(w http.ResponseWriter, r *http.Request){
	if r.Method == http.MethodPost{
		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil{
			//ADD REAL ERROR HERE!
			fmt.Println("FAILURE READING!!!")
			return
		}

		var receipt Receipt

		marshErr := json.Unmarshal(requestBody, &receipt)
		if marshErr != nil{
			fmt.Println("FAILED TO PARSE BECAUSE:", marshErr)
			return
		}

		fmt.Println("SUCCESS:", uuid.New().String())

		jsonString, toStrErr := json.Marshal(receipt)
		if toStrErr != nil {
			fmt.Println(toStrErr)
			return
		}

		fmt.Println("READ IN:\n", string(jsonString))

		//Calculate points here!

	}else{
		//ADD REAL ERROR HERE LATER!!!
		fmt.Println("BAD REQUEST!!!!")
	}
}

func main(){
	data := DataBase{map[string]int{}}

	http.HandleFunc("/receipts/process", data.InsertToDatabase)

	fmt.Println("Listening on Port 8000!")
	
	err := http.ListenAndServe(":8000", nil)

	if errors.Is(err, http.ErrServerClosed){
		fmt.Println("Server closed!")
	}else if err != nil{
		fmt.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}


}
