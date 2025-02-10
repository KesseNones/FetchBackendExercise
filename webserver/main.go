package main

import (
	"fmt"
	"net/http"
	"os"
	"errors"
	"github.com/google/uuid"
	"encoding/json"
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

type PointResponse struct {
	Points int `json:"points"`
}

func InvalidReceiptError(w http.ResponseWriter){
	http.Error(w, "The receipt is invalid", http.StatusBadRequest)
}

//Takes in the input receipt, calculates point value, 
// and writes points to database.
func (db *DataBase) InsertToDatabase(w http.ResponseWriter, r *http.Request){
	if r.Method == http.MethodPost{
		var receipt Receipt

		jsonParseErr := json.NewDecoder(r.Body).Decode(&receipt)
		if jsonParseErr != nil{
			InvalidReceiptError(w)
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

		//Kicks back error if any of the required fields don't exist 
		// or don't quite follow the desired format.
		if (
			len(receipt.Retailer) == 0 ||
			len(receipt.Total) == 0 ||
			receipt.Items == nil ||
			len(receipt.Items) == 0 ||
			len(strings.Split(receipt.PurchaseDate, "-")) != 3 ||
			len(strings.Split(receipt.PurchaseTime, ":")) != 2){
			
			InvalidReceiptError(w)
			return
		}

		prices := []float64{}

		//Iterates through all items in receipt, 
		// and makes sure required fields 
		// of each item have something in them, 
		// indicating they exist.
		//Also checks to make sure all prices 
		// are valid numbers and are positive.
		for _, item := range receipt.Items{
			if len(item.ShortDesc) == 0 || len(item.Price) == 0{
				InvalidReceiptError(w)
				return
			}
			//Kicks back error if any prices in items are non-numbers or negative.
			priceNum, err := strconv.ParseFloat(item.Price, 64)
			if err != nil || priceNum < 0.0{
				InvalidReceiptError(w)
				return
			}
			prices = append(prices, priceNum)
		}

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
			InvalidReceiptError(w)
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
		for i, item := range receipt.Items{
			if len(strings.TrimSpace(item.ShortDesc)) % 3 == 0{
				//Interpreting "number of points earned" to mean adding to total, 
				// not reseting total.
				totalPoints += int(math.Ceil(prices[i] * 0.2))
			}
		}

		//DEBUG; DESTROY LATER!!!
		fmt.Println(totalPoints)

		//Adds 6 points to total if the day of the purchase is odd.
		datePieces := strings.Split(receipt.PurchaseDate, "-")
		//Checks to make sure each part of date is a valid date piece.
		yearInt, yearErr := strconv.Atoi(datePieces[0])
		monthInt, monthErr := strconv.Atoi(datePieces[1])
		dayInt, dayErr := strconv.Atoi(datePieces[2])

		//Errors out if any part of the date isn't an integer.
		if yearErr != nil || monthErr != nil || dayErr != nil{
			InvalidReceiptError(w)
			return
		}

		//Checks to make sure month is within correct range.
		if monthInt < 1 || monthInt > 12{
			InvalidReceiptError(w)
			return
		}

		//Determines days of months for given year of receipt.
		daysOfMonths := []int{
			31, 28, 31, 
			30, 31, 30, 
			31, 31, 30, 
			31, 30, 31}
		//Makes February 29 days if it's a leap year.
		if yearInt % 4 == 0 && (yearInt % 100 != 0 || yearInt % 400 == 0){
			daysOfMonths[1] += 1
		}

		//Checks to make sure current day is a valid day in the given month.
		if dayInt < 1 || dayInt > daysOfMonths[monthInt - 1]{
			InvalidReceiptError(w)
			return
		}

		//Adds if day is odd.
		if dayInt % 2 == 1{
			totalPoints += 6
		}

		//DEBUG; DESTROY LATER!!!
		fmt.Println(totalPoints)

                //Adds 10 points if time of purchase is after 14:00 sharp and before 16:00.
                // This is interpreted to mean that any time between 14:01 
                // and 15:59 inclusive is valid.
		
		//Splits time into hours and minutes.
		timePieces := strings.Split(receipt.PurchaseTime, ":")
		
		//Parses hour and minute, kicking back error if either one isn't valid.
		hourInt, hourErr := strconv.Atoi(timePieces[0])
		minInt, minErr := strconv.Atoi(timePieces[1])
		if hourErr != nil || minErr != nil{
			InvalidReceiptError(w)
			return
		}

		//Throws error if hour and minute not in valid ranges for time.
		if hourInt < 0 || hourInt > 23 || minInt < 0 || minInt > 59{
			InvalidReceiptError(w)
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
		//Not likely to be thrown but here just in case.
		if idJsonErr != nil{
			http.Error(w, "Failed to encode ID JSON!", http.StatusInternalServerError)
			return
		}


	}else{
		//ADD REAL ERROR HERE LATER!!!
		fmt.Println("BAD REQUEST!!!!")
	}
}

//Given an id in the input url, accesses the database with that id, 
// and sends how many points that receipt was given.
func (db *DataBase) GetPointsFromId(w http.ResponseWriter, r *http.Request){
	//DEBUG; DESTROY LATER!
	fmt.Println(r.URL.String())

	if r.Method == http.MethodGet{
		//Fetches id to query database with.
		queryId := strings.Split(r.URL.String(), "/")[2]	

		//DEBUG; DESTROY LATER!
		fmt.Println(queryId)

		//Queries database using ID to get points associated with ID.
		points, found := db.Data[queryId]
		if found{
			response := PointResponse{points}
			w.Header().Set("Content-Type", "application/json")
			pointJsonErr := json.NewEncoder(w).Encode(response)
			if pointJsonErr != nil{
				//PUT REAL ERROR HERE LATER!
				fmt.Println("FAILED TO ENCODE POINTS FETCH JSON!")
				return
			}

		}else{
			//PUT REAL ERROR HERE LATER!
			fmt.Println("ITEM NOT FOUND!")
			return
		}

	}else{
		//ADD REAL ERROR HERE LATER!!!
		fmt.Println("INCORRECT REQUEST!! HAS TO BE GET!")
	}

}

func main(){
	data := DataBase{map[string]int{}}

	http.HandleFunc("/receipts/process", data.InsertToDatabase)
	http.HandleFunc("/receipts/{id}/points", data.GetPointsFromId)

	fmt.Println("Listening on Port 8000!")
	
	err := http.ListenAndServe(":8000", nil)

	if errors.Is(err, http.ErrServerClosed){
		fmt.Println("Server closed!")
	}else if err != nil{
		fmt.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}


}
