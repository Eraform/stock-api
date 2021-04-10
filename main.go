package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Stock struct {
	Ticker string  `json:"ticker"`
	Title  string  `json:"title"`
	Rsi    float32 `json:"rsi"`
}

var stocks []Stock

func getStocks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stocks)
}

func createStock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var stock Stock
	_ = json.NewDecoder(r.Body).Decode(&stock)
	stocks = append(stocks, stock)
	json.NewEncoder(w).Encode(&stock)
}
func getStock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range stocks {
		if item.Ticker == params["ticker"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Stock{})
}
func updateStock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range stocks {
		if item.Ticker == params["ticker"] {
			stocks = append(stocks[:index], stocks[index+1:]...)
			var Stock Stock
			_ = json.NewDecoder(r.Body).Decode(&Stock)
			Stock.Ticker = params["ticker"]
			stocks = append(stocks, Stock)
			json.NewEncoder(w).Encode(&Stock)
			return
		}
	}
	json.NewEncoder(w).Encode(stocks)
}
func deleteStock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range stocks {
		if item.Ticker == params["ticker"] {
			stocks = append(stocks[:index], stocks[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(stocks)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.Method, r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func main() {
	store, err := NewBadgerDB()
	if err != nil {
		log.Fatal(err)
	}
	defer store.db.Close()

	router := mux.NewRouter()

	stocks = append(stocks, Stock{Ticker: "MSFT", Title: "Microsoft Corp", Rsi: 69.05})
	stocks = append(stocks, Stock{Ticker: "DISCA", Title: "Discovery Inc.", Rsi: 33.75})

	router.HandleFunc("/stocks", getStocks).Methods("GET")
	router.HandleFunc("/stocks", createStock).Methods("Stock")
	router.HandleFunc("/stocks/{ticker}", getStock).Methods("GET")
	router.HandleFunc("/stocks/{ticker}", updateStock).Methods("PUT")
	router.HandleFunc("/stocks/{ticker}", deleteStock).Methods("DELETE")

	router.Use(loggingMiddleware)
	http.ListenAndServe(":8000", router)
}
