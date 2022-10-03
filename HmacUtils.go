package muxjwt

import (
	"crypto/hmac"
	"crypto/sha256"
	b64 "encoding/base64"
)

/*
func HmacEncodeStr encrypts a given string using HS256 by a given secret
@params:
	- str(string), the string to HS256 encrypt
	- secret(string), the secret to use for the encryption
@return: string, the HS256 encrypted string of str(param) via secret(param)
*/
func hmacEncodeStr(str string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(str))
	return b64.StdEncoding.EncodeToString(h.Sum(nil))
}

/*
func CmpHmacStr compares a HS256 string to a given string converted to it's HS256 form via a given secret
@params:
	- str(string), the string to encrypt to HS256 form
	- secret(string), the secret used to encode str(param) to it's HS256 form
	- hmac(string), the encrypted HS256 string to compare to
@return: bool, if the HS256 version of str(param) with secret(param) equals to a given HS256 string
*/
func cmpHmacStr(str string, secret string, hmac string) bool {
	return hmacEncodeStr(str, secret) == hmac
}
