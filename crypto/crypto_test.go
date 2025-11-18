package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"testing"
	"time"
)


func Test_Crypto_Hybrid_Envelope(t *testing.T){
	// Setup Keys (Simulated)
	privKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	pubKey := &privKey.PublicKey
	
	longSecret := "ThisSecretIsSuperLongAndWouldNormallyBreakRSAIfWeDidNotUseHybridEncryption_AABBCCDDEEFF"
	payload := fmt.Sprintf("%d|POST|%s|/api/v1/very/secure/resource", time.Now().Unix(), longSecret)

	secureToken, err := HybridEncrypt(pubKey, payload)

	if err != nil {
		t.Logf("Error: %s",err.Error())
	}
	t.Logf("Secure Token %s", secureToken)

	plainText, _ :=HybridDecrypt(privKey, secureToken)

	t.Logf("Secreto descifrado %s", plainText)


}