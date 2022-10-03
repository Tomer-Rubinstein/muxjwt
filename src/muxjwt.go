package main // debug
// package muxjwt

import (
	"encoding/json"
	"fmt"
	b64 "encoding/base64"
	"github.com/gorilla/mux"
	"crypto/hmac"
	"crypto/sha256"
	"net/http"
	"strings"
	"time"
	"errors"
)

type MuxJWT struct {
	r *mux.Router
	authenticate func(string, string) bool
	identify func() // used to check for token expiration, roles, etc..
}

type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type Payload struct {
	Sub string `json:"sub"`
	Iat int64 `json:"iat"`
}

func GenerateHeader() string {
	header := Header {
		Alg: "HS256",
		Typ: "JWT",
	}
	headerBytes, _ := json.Marshal(header)
	return b64.URLEncoding.EncodeToString(headerBytes)
}


func GeneratePayload(subject string, issuedAt int64) string {
	payload := Payload {
		Sub: subject,
		Iat: issuedAt,
	}
	payloadBytes, _ := json.Marshal(payload)
	return b64.URLEncoding.EncodeToString(payloadBytes)
}


func GenerateSignature(encodedHeader string, encodedPayload string, secret_salt string) string {
	return HmacEncodeStr(encodedHeader + "." + encodedPayload, secret_salt)
}


func HmacEncodeStr(str string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(str))
	return b64.StdEncoding.EncodeToString(h.Sum(nil))
}

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

	if !CmpHmacStr(token[0] + "." + token[1], token[2], secret) {
		return nil, errors.New("Invalid JWT format")
	}

	return payload, nil
}

func CmpHmacStr(str string, hmac string, secret string) bool {
	return HmacEncodeStr(str, secret) == hmac
}

func GenerateJWT(userid string, iat int64) string {
	encHeader := GenerateHeader()
	encPayload := GeneratePayload(userid, iat)
	encSignature := GenerateSignature(encHeader, encPayload, "DEBUG_SECRET")
	return encHeader + "." + encPayload + "." + encSignature
}

func NewMuxJWT(router *mux.Router, authFunc func(string, string) bool, identifyFunc func()) MuxJWT {
	router.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request){
		user := r.FormValue("user")
		password := r.FormValue("password")
		var jwt_token string

		if authFunc(user, password) == true {
			jwt_token = GenerateJWT(user, time.Now().Unix())
			fmt.Fprintf(w, jwt_token)
		}
	}).Methods("POST")

	return MuxJWT {
		r: router,
		authenticate: authFunc,
		identify: identifyFunc,
	}
}

func ProtectedRoute(r *mux.Router, route string, handler func(http.ResponseWriter, *http.Request)) *mux.Route {
	return r.HandleFunc(route, func(w http.ResponseWriter, r *http.Request){
		token := r.Header["Authorization"]
		if token == nil {
			fmt.Fprintf(w, "no auth token was given")
			return
		}

		_, err := TokenReadPayload(token[0], "DEBUG_SECRET", 60)
		if err != nil {
			fmt.Println(err)
			fmt.Fprintf(w, "Authentication error")
			return
		}

		// TODO: add roles

		fmt.Println(token)
		handler(w, r)
	})
}



/* DEBUG */
func main() {
  r := mux.NewRouter()
	NewMuxJWT(r, auth, identify)
	
  r.HandleFunc("/login", LoginHandler).Methods("GET")
	ProtectedRoute(r, "/secret", SecretHandler)
	fmt.Println("Listening on port 3000..")
	http.ListenAndServe(":3000", r)
}

func LoginHandler(w http.ResponseWriter, r *http.Request){
	http.ServeFile(w, r, "./static/LoginPage.html")
}

func SecretHandler(w http.ResponseWriter, r *http.Request){
	http.ServeFile(w, r, "./static/SecretPage.html")
}

func auth(username string, passw string) bool {
	return username == "admin" && passw == "admin"
}

func identify() { }
