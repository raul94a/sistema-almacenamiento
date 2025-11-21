package sdk

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)


func parseRSAPublicKeyFromPEM(pemStr string) (*rsa.PublicKey, error) {
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
