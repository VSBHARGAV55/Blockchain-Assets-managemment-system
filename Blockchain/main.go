package main

import (
	"github.com/BooMER23/Handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/createUser", Handlers.CreateUser).Methods("POST")
	router.HandleFunc("/sendMoney", Handlers.SendMoney).Methods("POST")
	router.HandleFunc("/receiveMoney", Handlers.ReceiveMoney).Methods("POST")
	router.HandleFunc("/printBlockchain", Handlers.PrintBlockchain).Methods("GET")
	router.HandleFunc("/Auth", Handlers.Authentication).Methods("GET")
	http.ListenAndServe(":8080", router)
}
