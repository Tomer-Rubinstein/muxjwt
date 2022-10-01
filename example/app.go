package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"example.com/muxjwt"
)

func main() {
  r := mux.NewRouter()
	mjwt := muxjwt.NewMuxJWT(r, auth, identify)
  muxjwt.protect(mjwt, r.HandleFunc("/login", LoginHandler).Methods("GET"))

	fmt.Println("Listening on port 3000..")
	http.ListenAndServe(":3000", r)
}

func LoginHandler(w http.ResponseWriter, r *http.Request){
	http.ServeFile(w, r, "./static/LoginPage.html")
}

func auth(username string, passw string) bool {
	return username == "admin" && passw == "admin"
}

func identify() { }
