package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"testing"
	"time"

	crypto_utils "github.com/storage-system/crypto/utils"
)

func Test_Crypto_Hybrid_Envelope(t *testing.T) {
	// Setup Keys (Simulated)
	privKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	pubKey := &privKey.PublicKey

	longSecret := "ThisSecretIsSuperLongAndWouldNormallyBreakRSAIfWeDidNotUseHybridEncryption_AABBCCDDEEFF"
	payload := fmt.Sprintf("%d|POST|%s|/api/v1/very/secure/resource", time.Now().Unix(), longSecret)

	secureToken, err := HybridEncrypt(pubKey, payload)

	if err != nil {
		t.Logf("Error: %s", err.Error())
	}
	t.Logf("Secure Token %s", secureToken)

	plainText, _ := HybridDecrypt(privKey, secureToken)

	t.Logf("Secreto descifrado %s", plainText)

}

func Test_Aes_Base_64_CipherText(t *testing.T) {
	k, _ := generateAesKey(AES_32_BYTES_KEY)
	t.Logf("Key: %s", string(k))
	payload := "TEST_AES_BASE_64_SPAIN_CIPHER"
	aesGcmCiphertextPacket, _ := aesGcmEncryption(k, payload)
	t.Log(aesGcmCiphertextPacket.GetBase64Ciphertext())
	t.Log(aesGcmCiphertextPacket.GetBase64Nonce())
	t.Logf("HashNonce: %s",aesGcmCiphertextPacket.GetNonceHash())
	t.Logf("Ciphertext: %s", string(aesGcmCiphertextPacket.Data))

	decrypted, _ := aesGcmDecryption(k, aesGcmCiphertextPacket)

	t.Logf("Decripted string: %s", decrypted)

	if payload != decrypted {
		t.Fatal("Ciphertext has not been decrypted correctly")
	}

}

func Test_Aes_Encryption_Decryption_Load(t *testing.T) {

	data := crypto_utils.CryptoUtilsTestData()
	for _, line := range data {
		k, _ := generateAesKey(AES_32_BYTES_KEY)
		aesGcmCiphertextPacket, _ := aesGcmEncryption(k, line)

		decrypted, _ := aesGcmDecryption(k, aesGcmCiphertextPacket)
		if line != decrypted {
			t.Logf("String to cipher: %s", line)
			t.Logf("Decrypted string: %s", decrypted)
			t.Fatalf("Ciphertext has not been decrypted correctly")

		}

	}
}


func Test_Pub_Key_To_String_Then_To_Pub_Key(t *testing.T){
	privKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	pubKey := &privKey.PublicKey
	s,e := RSAPublicKeyToPEM(pubKey)
	t.Logf("Pem: %s",s)
	if e != nil {
		t.Fatal(e.Error())
	}
	pk, e := ParseRSAPublicKeyFromPEM(s)
	if e != nil || (pk.N == nil){
		t.Fatal(e.Error())
	}
	if pk.E != pubKey.E {
		t.Fatal("Bad decoding of public key: bad exponent")	
	}
	res := pk.N.Cmp(pubKey.N)
	if pk.N.Cmp(pubKey.N) != 0 {
		t.Logf("Result of comparison pk vs pubKey %d",res)
		t.Logf("Decoded prime %v",pk.N)
		t.Logf("Original prime %v", pubKey.N)
		t.Fatal("Bad decoding of public key: bad prime number")	

	}
	t.Log("Encoding and decoding of RSA Public key completed successfully")
}