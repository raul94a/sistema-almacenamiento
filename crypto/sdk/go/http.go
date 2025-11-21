package sdk

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
)


func addRequestHeaders(req *http.Request, jwt *string, publicKey *rsa.PublicKey) (error){
	if jwt != nil {
		req.Header.Add("X-Authorization",*jwt)
	}
	site := req.Host
	req.Header.Add("X-Domain",site)
	utcString := time.Now().UTC().String()
	method := req.Method
	endpoint := req.URL
	payload := fmt.Sprintf("%s|%s|%s",utcString,method,endpoint)
	k,e := generateAesKey(AES_32_BYTES_KEY)
	if e != nil {
		return e
	}
	// The payload is ciphered with the AES Key
	packet, e := aesGcmEncryption(k,payload)
	if e != nil {
		return e
	}
	// The AES Key is then ciphered with the RSA public key
	encryptedAesKey, e := rsa.EncryptOAEP(sha256.New(),rand.Reader,publicKey,k,nil)
	if e != nil {
		return e
	}
	req.Header.Add("X-Auth",hex.EncodeToString(encryptedAesKey))
	req.Header.Add("X-Nonce", hex.EncodeToString(packet.Nonce))
	req.Header.Add("X-Digital-Envelope",hex.EncodeToString(packet.Data))

	return nil
	
}


