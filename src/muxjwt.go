package muxjwt

import (
	"encoding/json"
	"fmt"
	b64 "encoding/base64"
	"github.com/gorilla/mux"
	"crypto/hmac"
	"crypto/sha256"
	"net/http"
)

type MuxJWT struct {
	r *mux.Router
	authenticate func(string, string) bool
	identify func() // used to check for token expiration, roles, etc..
}


func GenerateHeader() string {
	type Header struct {
		Alg string `json:"alg"`
		Typ string `json:"typ"`
	}
	header := Header {
		Alg: "HS256",
		Typ: "JWT",
	}
	headerBytes, _ := json.Marshal(header)
	return b64.URLEncoding.EncodeToString(headerBytes)
}


func GeneratePayload(subject string, issuedAt int) string {
	type Payload struct {
		Sub string `json:"sub"`
		Iat int `json:"iat"`
	}
	payload := Payload {
		Sub: subject,
		Iat: issuedAt,
	}
	payloadBytes, _ := json.Marshal(payload)
	return b64.URLEncoding.EncodeToString(payloadBytes)
}


func GenerateSignature(encodedHeader string, encodedPayload string, secret_salt string) string {
	data := encodedHeader + "." + encodedPayload
	h := hmac.New(sha256.New, []byte(secret_salt))
	h.Write([]byte(data))

	return b64.StdEncoding.EncodeToString(h.Sum(nil))
}


func GenerateJWT(userid string, iat int) string {
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
			jwt_token = GenerateJWT(user, 999999999) // TODO: token expiration
			fmt.Fprintf(w, jwt_token)
		}
	}).Methods("POST")

	return MuxJWT {
		r: router,
		authenticate: authFunc,
		identify: identifyFunc,
	}
}

func protect(mjwt MuxJWT, r *mux.Route) {
	fmt.Println("works")
}
