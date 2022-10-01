package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"example.com/"
)

func main() {
  r := mux.NewRouter()
  r.HandleFunc("/login", LoginHandler).Methods("GET")
	r.HandleFunc("/auth", AuthHandler).Methods("POST")
	http.ListenAndServe(":80", r)
}

func LoginHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Welcome to the login page!")
}

func AuthHandler(w http.ResponseWriter, r *http.Request){
	// vars := mux.Vars(r)
  // w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "You've requested the book: %s on page %s\n", title, page)
}

