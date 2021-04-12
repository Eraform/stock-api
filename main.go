package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.Host, r.Method, r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func (s *Store) addStock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var stock Stock
	_ = json.NewDecoder(r.Body).Decode(&stock)
	err := s.Add(stock.Ticker, stock)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		json.NewEncoder(w).Encode(&stock)
	}
}

func (s *Store) getStocks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	stocks, err := s.GetStocks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(stocks)
}

func (s *Store) getStock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	stock, err := s.GetStock(params["ticker"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(stock)
}

func (s *Store) updateStock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	ticker := params["ticker"]

	var stock Stock
	_ = json.NewDecoder(r.Body).Decode(&stock)
	if ticker != stock.Ticker {
		http.Error(w, "params ticker and stock ticker do not match", http.StatusBadRequest)
		return
	}

	err := s.Upsert(ticker, stock)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(stock)
}

func (s *Store) deleteStock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	ticker := params["ticker"]

	err := s.DeleteItem(ticker)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	return
}

func (s *Store) getKeys(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	keys, err := s.GetAllKeys()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(keys)
}

func main() {
	path := getEnv("BADGER_PATH", "/tmp/badger-data")
	port := getEnv("PORT", "8080")

	store, err := NewBadger(path)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Badger().Close()

	router := mux.NewRouter()

	router.HandleFunc("/stocks", store.addStock).Methods(http.MethodPost)
	router.HandleFunc("/stocks", store.getStocks).Methods(http.MethodGet)
	router.HandleFunc("/stocks/{ticker}", store.getStock).Methods(http.MethodGet)
	router.HandleFunc("/stocks/{ticker}", store.updateStock).Methods(http.MethodPut)
	router.HandleFunc("/stocks/{ticker}", store.deleteStock).Methods(http.MethodDelete)
	router.HandleFunc("/keys", store.getKeys).Methods(http.MethodGet)

	router.Use(loggingMiddleware)
	http.ListenAndServe(":"+port, router)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
