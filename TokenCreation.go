package muxjwt

import (
	"encoding/json"
	b64 "encoding/base64"
	"time"
	"net/http"
	"fmt"
)


type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type Payload struct {
	Sub string `json:"sub"`
	Iat int64 `json:"iat"`
}

/*
func GenerateHeader converts the header part of JWT to a base64 string
@params: nil
@return: string, the base64 string of the JWT header part
*/
func generateHeader() string {
	header := Header {
		Alg: "HS256",
		Typ: "JWT",
	}
	headerBytes, _ := json.Marshal(header)
	return b64.URLEncoding.EncodeToString(headerBytes)
}

/*
func GeneratePayload sets 2 claims to the payload part of JWT: subject(str), issuedAt(int64)
as a Unix timestamp and returns it's base64 encoded form
@params:
	- subject(string), the subject claim(often used as userid or name)
@return: string, the base64 form of the payload part of JWT format
*/
func generatePayload(subject string) string {
	payload := Payload {
		Sub: subject,
		Iat: time.Now().Unix(),
	}
	payloadBytes, _ := json.Marshal(payload)
	return b64.URLEncoding.EncodeToString(payloadBytes)
}

/*
func GenerateSignature generates the signature part of JWT by concatinating the base64 versions
of the header and the payload by a dot and then HS256 encrypts the new strng using a given secret
@params:
	- encodedHeader(string), the base64 form of the header part of JWT
	- encodedPayload(string), the base64 form of the payload part of JWT
	- secret_salt(string), the secret to HS256 encrypt the described string with
@return: a HS256 encrypted string as described with secret_salt(param)
*/
func generateSignature(encodedHeader string, encodedPayload string, secret_salt string) string {
	return hmacEncodeStr(encodedHeader + "." + encodedPayload, secret_salt)
}

/*
func GenerateJWT creates a JWT string
@params:
	- userid(string), the subject claim to pass to the payload
@return: string, the JWT string: <base64(Header)>.<base64(Payload)>.<base64(Signature)>
*/
func (m MuxJWT) GenerateJWT(userid string) string {
	encHeader := generateHeader()
	encPayload := generatePayload(userid)
	encSignature := generateSignature(encHeader, encPayload, m.Secret)
	return encHeader + "." + encPayload + "." + encSignature
}

/*
func newCookie declares a new http.Cookie instance by given parameters
@params:
	- name(string), name of the cookie
	- value(string), value of the cookie
	- domain(string), to which domain this cookie points to
	- iat(int64), "issued at": the Unix timestamp at which the jwt was created
*/
func (m MuxJWT) NewCookie(name string, value string, domain string, iat int64) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = value
	cookie.Domain = domain
	cookie.RawExpires = fmt.Sprint(iat + m.ExpirationTime)
	cookie.Secure = false
	cookie.HttpOnly = true
	cookie.Raw = fmt.Sprintf("%s: %s", cookie.Name, cookie.Value)
	cookie.Unparsed = []string{cookie.Raw}

	return cookie
}
