package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/paladinknightmaster/currency-exchange-go-backend/handlers"
)

var godotenvLoadFunc = godotenv.Load
var serverStartFunc = http.ListenAndServe

func init() {
	loadEnv()
}

func loadEnv() {
	if err := godotenvLoadFunc(); err != nil {
		log.Println("No .env file found")
	}
}

func main() {
	if err := startServer(); err != nil {
		log.Fatal(err)
	}
}

func startServer() error {
	r := setupRouter()
	log.Println("Server is running on port 8080")
	err := serverStartFunc(":8080", r)
	if err != nil {
		return err
	}
	return nil
}

func setupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/rates", handlers.GetRates).Methods("GET")
	r.HandleFunc("/convert", handlers.ConvertCurrency).Methods("GET")
	return r
}
