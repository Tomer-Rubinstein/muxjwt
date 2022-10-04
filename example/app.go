package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/Tomer-Rubinstein/muxjwt"
	"fmt"
)

func main() {
  r := mux.NewRouter()
	muxjwt.InitAuthRoute(r, authFunc, "/auth", identifyFunc, "user", "password")

  r.HandleFunc("/login", LoginHandler).Methods("GET")
	muxjwt.ProtectedRoute(r, "/secret", SecretHandler).Methods("GET")
	fmt.Println("Listening on port 3000..")
	http.ListenAndServe(":3000", r)
}

func LoginHandler(w http.ResponseWriter, r *http.Request){
	http.ServeFile(w, r, "./static/LoginPage.html")
}

func SecretHandler(w http.ResponseWriter, r *http.Request){
	http.ServeFile(w, r, "./static/SecretPage.html")
}

func authFunc(body map[string]string) bool {
	username := body["user"]
	password := body["password"]
	fmt.Println(username, password)
	return username == "admin" && password == "admin"
}

func identifyFunc(body map[string]string) string {
	return body["user"]
}
