package cryptoutils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type KeyBearer[T any] interface {
	GetKey(url string) (T,error)
}

type StorageSystemRsaPublicKey struct {
	PublicKey *rsa.PublicKey
	Url string
}

func (s *StorageSystemRsaPublicKey) GetKey()(*rsa.PublicKey,error){
	req, err := http.NewRequest("GET",s.Url, nil)
	if err != nil {
		return nil,err
	}
	client := &http.Client{
		Timeout: time.Second * 30, // Buen hábito en SDKs
	}

	res, err := client.Do(req)
	if err != nil {
		return nil,err
	}

	var pkey struct {
		PublicKey string `json:"public_key"`
	}
	body := res.Body
	defer body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error API: status %d", res.StatusCode)
	}
	
	if err := json.NewDecoder(res.Body).Decode(&pkey); err != nil {
		return nil, err
	}

	// Now we have the PEM that the server will give us so now it's time to transform it to rsa.PublicKey
	pubKey, err := s.parsePemToRsaPublicKey(pkey.PublicKey)
	if err != nil {
		return nil,err
	}

	return pubKey,nil
	

}


func (s *StorageSystemRsaPublicKey) parsePemToRsaPublicKey(pemStr string) (*rsa.PublicKey, error) {
    block, _ := pem.Decode([]byte(pemStr))
    if block == nil {
        return nil, errors.New("failed to decode PEM block")
    }

    var pub *rsa.PublicKey
    var err error

    // Try PKCS1 first (traditional "RSA PUBLIC KEY")
    if block.Type == "RSA PUBLIC KEY" || block.Type == "PUBLIC KEY" {
        pub, err = x509.ParsePKCS1PublicKey(block.Bytes)
        if err == nil {
            return pub, nil
        }
    }

    // If that failed, try PKIX ("PUBLIC KEY" header, works for any key type)
    if block.Type == "PUBLIC KEY" {
        parsedKey, err2 := x509.ParsePKIXPublicKey(block.Bytes)
        if err2 != nil {
            return nil, err2 // or combine both errors if you want
        }
        var ok bool
        pub, ok = parsedKey.(*rsa.PublicKey)
        if !ok {
            return nil, errors.New("not an RSA public key")
        }
        return pub, nil
    }

    return nil, errors.New("unsupported PEM type: " + block.Type)
}
