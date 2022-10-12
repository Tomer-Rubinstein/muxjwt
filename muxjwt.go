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
// TODO: add proper logs

type MuxJWT struct {
	Secret string
	ExpirationTime int64
	Host string
}

/*
func NewMuxJWT returns a new MuxJWT instance
@params:
	- secret(string), the secret to use for the JWT encryption
	- expTime(int64), the life-time in seconds of every JWT
	- host(string), the host address
@return: MuxJWT
*/
func NewMuxJWT(secret string, expTime int64, host string) MuxJWT {
	if secret == "" {
		panic("MuxJWT: secret must not be an empty string")
	}
	if expTime <= 0 {
		panic("MuxJWT: expiration time shouldn't be nil(0) nor negative")
	}
	if host == "" {
		panic("MuxJWT: host undefined")
	}

	return MuxJWT {
		Secret: secret,
		ExpirationTime: expTime,
		Host: host,
	}
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
func (m MuxJWT) InitAuthRoute(router *mux.Router, authFunc func(map[string]string) bool, authRoute string, identifyFunc func(map[string]string)string, bodyKeys ...string) {
	router.HandleFunc(authRoute, func(w http.ResponseWriter, r *http.Request){
		// TODO?: accept JSON as POST data
		body := make(map[string]string)
		for i:=0; i < len(bodyKeys); i++ {
			body[bodyKeys[i]] = r.FormValue(bodyKeys[i])
		}

		var jwt_token string
		fmt.Println(body)
		if authFunc(body) {
			jwt_token = m.GenerateJWT(identifyFunc(body))
			jwt_payload, err := m.TokenReadPayload(jwt_token)
			if err != nil {
				panic("Invalid JWT: couldn't read payload")
			}
			cookie := m.NewCookie("token_"+m.Host, jwt_token, m.Host, jwt_payload.(Payload).Iat)

			http.SetCookie(w, cookie)
			w.WriteHeader(http.StatusAccepted) // successful auth
			return
		}
		w.WriteHeader(http.StatusForbidden) // failed auth
	}).Methods("POST")
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
func (m MuxJWT) TokenReadPayload(jwt string) (interface{}, error) {
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

	if time.Now().Unix() - payload.Iat >= m.ExpirationTime {
		return nil, errors.New("Token is expired")
	}

	if !cmpHmacStr(token[0] + "." + token[1], m.Secret, token[2]) {
		return nil, errors.New("Invalid JWT format")
	}

	return payload, nil
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
func (m MuxJWT) ProtectedRoute(r *mux.Router, route string, handler func(http.ResponseWriter, *http.Request)) *mux.Route {
	return r.HandleFunc(route, func(w http.ResponseWriter, r *http.Request){
		tokenCookie, err := r.Cookie("token_"+m.Host)
		if err != nil {
			m.NoAuthPage(w)
			return
		}

		_, err = m.TokenReadPayload(tokenCookie.Value)
		if err != nil {
			m.NoAuthPage(w)
			return
		}
		// TODO: roles
		handler(w, r)
	})
}

/*
func DeleteJWTCookie deletes a muxjwt cookie from the client's webbrowser
@params:
	- w(http.ResponseWriter), response writer
@return: nil
*/
func (m MuxJWT) DeleteJWTCookie(w http.ResponseWriter) {
	c := m.NewCookie("token_"+m.Host, "", m.Host, -m.ExpirationTime)
	c.MaxAge = -1
	http.SetCookie(w, c)
}

func (m MuxJWT) NoAuthPage(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte("No authentication"))
}
