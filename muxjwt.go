package muxjwt

import (
	"encoding/json"
	"fmt"
	b64 "encoding/base64"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
	"errors"
	"time"
)

/*
func TokenReadPayload validates a given jwt (including expiration) and reads the
payload of the jwt.
@params:
	- jwt(string), the JSON Web Token
	- secret(string), the secret salt used to encode the token
	- expirationSec(int64), the life-time of the token in seconds
@return: (
	- interface{}, nil if error occurred, otherwise, parsed Payload struct
	- error
)
*/
func TokenReadPayload(jwt string, secret string, expirationSec int64) (interface{}, error) {
	token := strings.Split(jwt, ".") // '.' isn't a base64 character
	if len(token) != 3 {
		return nil, errors.New("Token isn't consisted of 3 parts: HEADER.PAYLOAD.SIGNATURE")
	}

	decodedPayload, err := b64.URLEncoding.DecodeString(token[1])
	if err != nil {
		return nil, errors.New("PAYLOAD isn't a valid base64 string")
	}

	payload := Payload{}
	err = json.Unmarshal([]byte(decodedPayload), &payload)
	if err != nil {
		return nil, err
	}

	if time.Now().Unix() - payload.Iat >= expirationSec {
		return nil, errors.New("Token is expired")
	}

	if !CmpHmacStr(token[0] + "." + token[1], secret, token[2]) {
		return nil, errors.New("Invalid JWT format")
	}

	return payload, nil
}

/*
func InitAuthRoute initializes a new auth route using Gorilla Mux.
this route is used for authenticating user credentials given by POST request body parameters
and validates the credentials based on a given auth function.
@params:
	- router(*mux.Router), the router instance of the app
	- authFunc(func(...string)->bool), the authentication function to validate given user credentials
	- authRoute(string), the route to bind the auth service to
	- bodyKeys(...string), the names(keys) of the POST body parameters to give to authFunc as values (ORDER MATTERS!)
@return: nil (void)
*/
func InitAuthRoute(router *mux.Router, authFunc func(map[string]string) bool, authRoute string, bodyKeys ...string) {
	router.HandleFunc(authRoute, func(w http.ResponseWriter, r *http.Request){
		// TODO?: accept JSON as POST data
		body := make(map[string]string)
		for i:=0; i < len(bodyKeys); i++ {
			body[bodyKeys[i]] = r.FormValue(bodyKeys[i])
		}

		var jwt_token string
		if authFunc(body) == true {
			jwt_token = GenerateJWT(body["username"])
			fmt.Fprintf(w, jwt_token) // TODO: use gorilla/securecookie
		}
	}).Methods("POST")
}

/*
func ProtectedRoute creates a new route using Gorilla Mux that can only be accessed to by using a valid JWT
in the request header "Authorization"
@params:
	- r(*mux.Router), the router instance of the app
	- route(string), the route to protect
	- handler(func(http.ResponseWriter, *http.Request)->Any), the handler function to the route
@return: *mux.Route, so you can continue using this route as a normal r.HandleFunc() struct type
*/
func ProtectedRoute(r *mux.Router, route string, handler func(http.ResponseWriter, *http.Request)) *mux.Route {
	return r.HandleFunc(route, func(w http.ResponseWriter, r *http.Request){
		// TODO: what is the need of the "Bearer" prefix?
		token := strings.Trim(r.Header["Authorization"][0], "Bearer ")
		if token == "" {
			fmt.Fprintf(w, "No auth token was given")
			return
		}

		_, err := TokenReadPayload(token, "DEBUG_SECRET", 60) // TODO: config
		if err != nil {
			fmt.Println(err)
			fmt.Fprintf(w, "Authentication error")
			return
		}
		// TODO: add roles
		handler(w, r)
	})
}








/* DEBUG */
func main() {
  r := mux.NewRouter()
	InitAuthRoute(r, authFunc, "/auth", "username", "password")
	
  r.HandleFunc("/login", LoginHandler).Methods("GET")
	ProtectedRoute(r, "/secret", SecretHandler).Methods("GET")
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
	username := body["username"]
	password := body["password"]
	return username == "admin" && password == "admin"
}