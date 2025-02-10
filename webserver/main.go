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
	"strings"
	"math"
)

type ReceiptItem struct{
	ShortDesc string `json:"shortDescription"`
	Price string `json:"price"`
}

//Input struct used to handle either receipt input or just id.
type Receipt struct{
	Retailer string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Total string `json:"total"`
	Items []ReceiptItem `json:"items"`
}

type DataBase struct {
	Data map[string]int
}

type IdResponse struct {
	Id string `json:"id"`
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

		//5 points for every two items on receipt.
		totalPoints += ((len(receipt.Items) / 2) * 5)

		//DEBUG; DESTROY LATER!!!
		fmt.Println(totalPoints)

		//Iterates through items of given receipt, 
		// if the trimmed length of an item is a multiple of 3,
		// multiplies the item price by 0.2 
		// and rounds up to nearest integer (ceiling), 
		// the result of which is the number of points added 
		// to the total from that given item.
		for _, item := range receipt.Items{
			if len(strings.TrimSpace(item.ShortDesc)) % 3 == 0{
				floatPrice, err := strconv.ParseFloat(item.Price, 64)
				if err != nil{
					//PUT REAL ERROR HERE LATER!!!
					fmt.Println("FAILED TO PARSE ITEM PRICE!")
					return
				}

				//Interpreting "number of points earned" to mean adding to total, 
				// not reseting total.
				totalPoints += int(math.Ceil(floatPrice * 0.2))
			}
		}

		//DEBUG; DESTROY LATER!!!
		fmt.Println(totalPoints)

		//Adds 6 points to total if the day of the purchase is odd.
		datePieces := strings.Split(receipt.PurchaseDate, "-")
		dayInt, dayErr := strconv.Atoi(datePieces[2])
		if dayErr != nil{
			//PUT REAL ERROR HERE LATER!
			fmt.Println("FAILED TO PARSE DAY!")
			return
		}
		if dayInt % 2 == 1{
			totalPoints += 6
		}

		//DEBUG; DESTROY LATER!!!
		fmt.Println(totalPoints)

                //Adds 10 points if time of purchase is after 14:00 sharp and before 16:00.
                // This is interpreted to mean that any time between 14:01 
                // and 15:59 inclusive is valid.
		// Also assuming that hour stays in range 0-23 
		// inclusive and minute 0-59 inclusive.
		//Splits time into hours and minutes.
		timePieces := strings.Split(receipt.PurchaseTime, ":")
		//Parses hour.
		hourInt, hourErr := strconv.Atoi(timePieces[0])
		if hourErr != nil{
			//PUT REAL ERROR HERE LATER!
			fmt.Println("ERROR! Failed to parse purchase hour!")
			return
		}
		//Parses minute.
		minInt, minErr := strconv.Atoi(timePieces[1])
		if minErr != nil{
			//PUT REAL ERROR HERE LATER!
			fmt.Println("ERROR! Failed to parse purchase minute!")
			return
		}

		//Adds 10 points if within desired range.
		if (hourInt == 14 && minInt > 0) || hourInt == 15{
			totalPoints += 10
		} 


		//DEBUG; DESTROY LATER!!!
		fmt.Println(totalPoints)

		//Generates uuid for receipt and inserts into database.
		id := uuid.New().String()
		db.Data[id] = totalPoints
		
		output := IdResponse{id}

		//Sends response json with id.
		w.Header().Set("Content-Type", "application/json")
		idJsonErr := json.NewEncoder(w).Encode(output)
		if idJsonErr != nil{
			//PUT REAL ERROR HERE LATER!
			fmt.Println("FAILED TO ENCODE RESPONSE JSON!")
			return
		}


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
