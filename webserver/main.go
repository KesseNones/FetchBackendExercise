package main

import (
	"fmt"
	"net/http"
	"os"
	"errors"
	"github.com/google/uuid"
	"encoding/json"
	"io/ioutil"
	"unicode"
	"strconv"
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
			//PUT REAL ERROR HERE LATER!!!
			fmt.Println("FAILED TO PARSE BECAUSE:", marshErr)
			return
		}

		//DEBUG; DESTROY LATER
		fmt.Println("SUCCESS:", uuid.New().String())

		//DEBUG OUTPUT
		//DESTROY LATER!!!
		jsonString, toStrErr := json.Marshal(receipt)
		if toStrErr != nil {
			fmt.Println(toStrErr)
			return
		}
		fmt.Println("READ IN:\n", string(jsonString))

		totalPoints := 0

		//One point for every alphanumeric character 
		// in retailer name.
		for _, c := range receipt.Retailer{
			if unicode.IsLetter(c) || unicode.IsDigit(c){
				totalPoints += 1
			} 
		}

		//DEBUG; DESTROY LATER
		fmt.Println(totalPoints)

		//Parses receipt Total to floating point 
		// for the two following point calculations.
		floatTotal, convErr := strconv.ParseFloat(receipt.Total, 64)
		if convErr != nil{
			//PROBABLY ADD A REAL ERROR HERE LATER
			fmt.Println("FAILED TO CONVERT TOTAL TO FLOAT!!!")
			return
		}

		//50 points if total is round dollar amount.
		//Might have edge case with fractional cents, 
		// if testing goes that far.
		if floatTotal == float64(int(floatTotal)){
			totalPoints += 50
		}

		//DEBUG; DESTROY LATER!!!
		fmt.Println(totalPoints)

		//25 points if the total is a multiple of 0.25
		if int(floatTotal * 100) % 25 == 0{
			totalPoints += 25
		}
		//DEBUG; DESTROY LATER!!!
		fmt.Println(totalPoints)

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
