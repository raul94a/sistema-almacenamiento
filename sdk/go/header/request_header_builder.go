package request_header

import (
	"fmt"
	"net/http"
)

type header string
const (
	X_AUTHORIZATION header = "X-Authorization"
	X_NONCE header= "X-Nonce"
	X_SITE header= "X-Domain"
	X_DIGITAL_ENVELOPE header= "X-Digital-Envelope"
	X_AUTH header= "X-Auth"
	X_HASH header= "X-Hash"
)


type RequestHeaderBuilder struct {
	Request *http.Request

	
}

func (r *RequestHeaderBuilder) Builder(req *http.Request) *RequestHeaderBuilder {
	r.Request = req
	return r
}

func (r *RequestHeaderBuilder) AddJwtHeader(jwt string) *RequestHeaderBuilder {
	if r.Request == nil {
		panic(fmt.Errorf("RequestHeaderBuilder error: Call .Builder(req) before attempting to add headers"))
	}
	r.Request.Header.Add(string(X_AUTHORIZATION),jwt)
	return r
}

func (r *RequestHeaderBuilder) AddSite() *RequestHeaderBuilder {
	if r.Request == nil {
		panic(fmt.Errorf("RequestHeaderBuilder error: Call .Builder(req) before attempting to add headers"))
	}
	r.Request.Header.Add(string(X_SITE),r.Request.Host)
	return r
}

// `digitalEnvelope` is the hex-encoded ciphertext
func (r *RequestHeaderBuilder) AddDigitalEnvelope(digitalEnvelope string) *RequestHeaderBuilder {
	if r.Request == nil {
		panic(fmt.Errorf("RequestHeaderBuilder error: Call .Builder(req) before attempting to add headers"))
	}
	r.Request.Header.Add(string(X_DIGITAL_ENVELOPE),digitalEnvelope)
	return r
}
// `auth` is the hex-encoded encrypted AES Key with the RSA Public Key
// The RSA Public Key is provided by the Storage System Server
// The encryption of the AES Key with the RSA Public Key is a responsability
// of the cryptoutils package.
func (r *RequestHeaderBuilder) AddAuth(auth string) *RequestHeaderBuilder {
	if r.Request == nil {
		panic(fmt.Errorf("RequestHeaderBuilder error: Call .Builder(req) before attempting to add headers"))
	}
	r.Request.Header.Add(string(X_AUTH),auth)
	return r
}
// `nonce` is the hex-encoded nonce used in the AES Encryption of the payload
func (r *RequestHeaderBuilder) AddNonce(nonce string) *RequestHeaderBuilder {
	if r.Request == nil {
		panic(fmt.Errorf("RequestHeaderBuilder error: Call .Builder(req) before attempting to add headers"))
	}
	r.Request.Header.Add(string(X_NONCE),nonce)
	return r
}
// Nonce SHA-3 256 Bytes length hash
func (r *RequestHeaderBuilder) AddHashHeader(hash string) *RequestHeaderBuilder {
	if r.Request == nil {
		panic(fmt.Errorf("RequestHeaderBuilder error: Call .Builder(req) before attempting to add headers"))
	}
	r.Request.Header.Add(string(X_HASH),hash)
	return r
}

func (r *RequestHeaderBuilder) Build() *http.Request{
	if r.Request == nil {
		panic(fmt.Errorf("RequestHeaderBuilder error: Call .Builder(req) before attempting to add headers"))
	}
	return r.Request
}