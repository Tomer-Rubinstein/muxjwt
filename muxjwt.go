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

// TODO: implement gorilla/securecookie(s) instead
// TODO: implement logout functionality

var Secret string
var ExpirationTime int64 // Note: in seconds
var Host string

func init(){
	Secret = "DEBUG_SECRET"
	ExpirationTime = 60
	Host = "localhost"
}

/*
func TokenReadPayload validates a given jwt (including expiration) and reads the
payload of the jwt.
@params:
	- jwt(string), the JSON Web Token
@return: (
	- interface{}, nil if error occurred, otherwise, parsed Payload struct
	- error
)
*/
func TokenReadPayload(jwt string) (interface{}, error) {
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

	if time.Now().Unix() - payload.Iat >= ExpirationTime {
		return nil, errors.New("Token is expired")
	}

	if !cmpHmacStr(token[0] + "." + token[1], Secret, token[2]) {
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
func InitAuthRoute(router *mux.Router, authFunc func(map[string]string) bool, authRoute string, identifyFunc func(map[string]string)string, bodyKeys ...string) {
	router.HandleFunc(authRoute, func(w http.ResponseWriter, r *http.Request){
		// TODO?: accept JSON as POST data
		body := make(map[string]string)
		for i:=0; i < len(bodyKeys); i++ {
			body[bodyKeys[i]] = r.FormValue(bodyKeys[i])
		}

		var jwt_token string
		fmt.Println(body)
		if authFunc(body) == true {
			jwt_token = generateJWT(identifyFunc(body))
			jwt_payload, err := TokenReadPayload(jwt_token)
			if err != nil {
				panic("Invalid JWT: couldn't read payload")
			}
			cookie := newCookie("token_"+Host, jwt_token, Host, "/secret", jwt_payload.(Payload).Iat)

			http.SetCookie(w, cookie)
			fmt.Fprintf(w, jwt_token) // DEBUG
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
		tokenCookie, err := r.Cookie("token_"+Host)
		if err != nil {
			fmt.Printf("Error occured while reading testcookie")
			return
		}

		_, err = TokenReadPayload(tokenCookie.Value)
		if err != nil {
			fmt.Println(err)
			fmt.Fprintf(w, "Authentication error") // TODO: add customizability
			return
		}
		// TODO: roles
		handler(w, r)
	})
}
