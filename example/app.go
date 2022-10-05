package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/Tomer-Rubinstein/muxjwt"
	"fmt"
)

func main() {
  r := mux.NewRouter()
	m := muxjwt.NewMuxJWT("my_secret", 60*60*24*7, "localhost")
	m.InitAuthRoute(r, authFunc, "/auth", identifyFunc, "username", "password")

  r.HandleFunc("/login", LoginHandler).Methods("GET")
	m.ProtectedRoute(r, "/secret", SecretHandler).Methods("GET")

	fmt.Println("Listening on port 3000..")
	http.ListenAndServe(":3000", r)
}

func LoginHandler(w http.ResponseWriter, r *http.Request){
	http.ServeFile(w, r, "./static/LoginPage.html")
}

func SecretHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "<h1>Welcome to the Secret Page!</h1>")
}

func authFunc(body map[string]string) bool {
	username := body["username"]
	password := body["password"]
	fmt.Println(username, password)
	return username == "admin" && password == "admin"
}

func identifyFunc(body map[string]string) string {
	return body["username"]
}
