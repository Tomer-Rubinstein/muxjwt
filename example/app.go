package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"example.com/muxjwt"
)

func main() {
  r := mux.NewRouter()
	jwt := muxjwt.NewMuxJWT(r, auth, identify)
  r.HandleFunc("/login", LoginHandler).Methods("GET")
	http.ListenAndServe(":80", r)
}

func LoginHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Welcome to the login page!")
}

func auth(username string, passw string) bool {
	if username == "admin" && passw == "admin" {
		return true
	}
	return false
}

func identify() { }
