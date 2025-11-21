package crypto

/*

	For keeping the app secure, we're providing 3 layers of security:

	A) JWT Authorization.
	B) Hybrid RSA/AES Encryption
	C) Domain whitelist


	A) JWT Authorization is not directly implemented by the storage system.
	   it relies on a third-party that have to expose its JWKS.
	   With these JWKS the JWT Token can be decyphered and used to
	   Authenticate the user. Special fields will be allowed to be included
	   into the JWT through an admin user interface.

	B) Hybrid RSA/AES Encryption. Provided by the Storage System. This will
	   be used as a method to verify the requests. The final output of
	   this technique will be a encrypted string, the `Signature`, included
	   in the headers. The secrets and keys can be created/modified through the
	   admin user interface. Also, some endpoints will be provided for automation
	   of key rotations.

	C) Domain whitelist. Registration of the domains that can talk to our system.
	   If the user lack of firewall knowledge this will be pretty useful.
	   The Storage system is able to accept domains that are allowed to talk
	   to the Storage System.

	All of the above security layers can be used in any combination of them. Also,
	any of these layers can be disabled through the admin user interface.

	Needs to this system can be successful:

	a) SDKs
	   We do need to build SDKs that manages the request to the storage system.
	   These will manage all the complex encryption needed for the Signature.

	   The endpoints will be available with simple methods.


	   # Kotlin


	   ```kotlin
	   // @param key: Storage System Key
	   // @param secret: Storage System Secret
	   data class StorageSystemConfig(private val key: String, private) {

	   }
	   public class StorageSystemClient(val config: StorageSystemConfig? = null) {



	   }
	   ```


*/
import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha3"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"
)

type aesKeySize int

const (
	AES_16_BYTES_KEY aesKeySize = 16
	AES_24_BYTES_KEY aesKeySize = 24
	AES_32_BYTES_KEY aesKeySize = 32
)

// private function used to generate an AES Key.
// The AES Key is 16, 24 or 32 bytes long.
func generateAesKey(size aesKeySize) ([]byte, error) {
	// should we keep the session key safe?
	aesKey := make([]byte, size)
	length, err := io.ReadFull(rand.Reader, aesKey)
	if err != nil {
		// Handle the error, as the buffer was not completely filled.
		// For a cryptographic operation like this, a non-nil error is a critical failure.
		return nil, fmt.Errorf("error reading random bytes: %w", err)
	}
	if aesKeySize(length) != size {
		return nil, fmt.Errorf("read length does not match size")
	}
	return aesKey, nil
}

type AesGcmCiphertextPacket struct {
	Nonce []byte
	Data  []byte
}

// Transform the Nonce bytes in a base64 String
// The `Nonce` is needed to encrypt and decrypt the ciphertext.
// The base64 String can be used to retrieve the String Nonce.
// The Nonce must not be exposed to any client. AES/GCM is a
// symmetric encryption algorithm. By Knowing the key and the nonce,
// the system is compromised.
func (packet *AesGcmCiphertextPacket) GetBase64Nonce() string {
	return base64.StdEncoding.EncodeToString(packet.Nonce)
}

// GetNonceHash calculates the SHA3-256 hash of the packet's nonce.
// 
// This hash serves two primary purposes:
// 1. Integrity Check: It verifies that the nonce value included in the
//    request has not been tampered with in transit (data integrity).
// 2. Binding to Signature: This hash is included in the data that is
//    digitally signed by the sender's private key. This ensures the 
//    signature is unique for this specific request, which is crucial for 
//    preventing **replay attacks**.
// 
// The output hash is always 256 bits (32 bytes) and is returned 
// as a 64-character hexadecimal string.
// 
// NOTE: **Authenticity** (verifying the sender's identity) is handled
// by successfully decrypting the digital signature using the sender's
// public key, not by comparing the nonce hash alone.
func (packet *AesGcmCiphertextPacket) GetNonceHash() string {
	h := sha3.New256()
	h.Write(packet.Nonce)
	digest := h.Sum(nil)
	hexEncoding := hex.EncodeToString(digest)
	return hexEncoding
}

// Transform the Ciphertext bytes into a base64 String
func (packet *AesGcmCiphertextPacket) GetBase64Ciphertext() string {
	return base64.StdEncoding.EncodeToString(packet.Data)
}

// Use the AES Key generated in `generateAesKey` to cipher
// a message.
func aesGcmEncryption(key []byte, payload string) (*AesGcmCiphertextPacket, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	// Seal encrypts and authenticates the data
	cipherData := gcm.Seal(nil, nonce, []byte(payload), nil)
	aesGcmCiphertextPacket := &AesGcmCiphertextPacket{
		Nonce: nonce,
		Data:  cipherData,
	}
	return aesGcmCiphertextPacket, nil
}

func aesGcmDecryption(key []byte, packet *AesGcmCiphertextPacket) (string, error) {
	const NONCE_SIZE = 12
	// C. Decrypt the Payload using the recovered AES Session Key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	if len(packet.Nonce) != NONCE_SIZE {
		// Special error! watch out
		return "", fmt.Errorf("Error f117")
	}

	plaintextBytes, err := gcm.Open(nil, packet.Nonce, packet.Data, nil)
	if err != nil {
		// Special error! watch out
		return "", fmt.Errorf("Error f118")
	}
	return string(plaintextBytes), nil
}


func RSAPublicKeyToPEM(pub *rsa.PublicKey) (string, error) {
    // Marshal to PKCS1 (traditional format)
    derBytes := x509.MarshalPKCS1PublicKey(pub)

    // Or use PKIX if you prefer the newer standard (most tools accept both)
    // derBytes, err := x509.MarshalPKIXPublicKey(pub)
    // if err != nil { return "", err }

    block := &pem.Block{
        Type:  "RSA PUBLIC KEY", // or "PUBLIC KEY" if using PKIX
        Bytes: derBytes,
    }

    pemBytes := pem.EncodeToMemory(block)
    return string(pemBytes), nil
}
func ParseRSAPublicKeyFromPEM(pemStr string) (*rsa.PublicKey, error) {
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


// SecureEnvelope holds the two parts of our hybrid encryption
type SecureEnvelope struct {
	EncryptedKey  string `json:"k"` // RSA-encrypted AES key (Base64)
	EncryptedData string `json:"d"` // AES-encrypted Payload (Base64)
	Nonce         string `json:"n"` // AES Nonce/IV (Base64)
}

// --- 1. HYBRID ENCRYPT (Client Side) ---
func HybridEncrypt(pubKey *rsa.PublicKey, payload string) (string, error) {
	// Security errors
	//*
	// 1. Check payload size. Put a max
	// 2. chech nonce hgas 12 bytes after decode
	// 3. CHECK ERRORS IN DECODE
	// 4. CHECK PUB KEY LENGTH!! SHOULD BE >=2048 AND POWER OF 2
	// 5. Add nonce hash to  PAYLOAD this will help to check against attacks
	// *//
	// A. Generate a random 32-byte AES key (The "Session Key")
	sessionKey := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, sessionKey); err != nil {
		return "", err
	}

	// B. Encrypt the Payload using AES-GCM (Authentication + Encryption)
	block, err := aes.NewCipher(sessionKey)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	// Seal encrypts and authenticates the data
	cipherData := gcm.Seal(nil, nonce, []byte(payload), nil)

	// C. Encrypt the Session Key using RSA (The "Digital Envelope")
	encryptedKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pubKey, sessionKey, nil)
	if err != nil {
		return "", err
	}

	// D. Bundle everything into a JSON string
	packet := SecureEnvelope{
		EncryptedKey:  base64.StdEncoding.EncodeToString(encryptedKey),
		EncryptedData: base64.StdEncoding.EncodeToString(cipherData),
		Nonce:         base64.StdEncoding.EncodeToString(nonce),
	}

	jsonBytes, _ := json.Marshal(packet)
	return string(jsonBytes), nil
}

// --- 2. HYBRID DECRYPT (Server Side) ---
func HybridDecrypt(privKey *rsa.PrivateKey, jsonPacket string) (string, error) {
	var packet SecureEnvelope
	if err := json.Unmarshal([]byte(jsonPacket), &packet); err != nil {
		return "", fmt.Errorf("bad json format")
	}

	// A. Decode Base64 parts
	// MUST CHECK ERROR
	encKeyBytes, _ := base64.StdEncoding.DecodeString(packet.EncryptedKey)
	encDataBytes, _ := base64.StdEncoding.DecodeString(packet.EncryptedData)
	nonceBytes, _ := base64.StdEncoding.DecodeString(packet.Nonce) // decode must have 12 bytes

	// B. Decrypt the Session Key using RSA Private Key
	sessionKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privKey, encKeyBytes, nil)
	if err != nil {
		return "", fmt.Errorf("RSA decryption failed (Wrong Private Key?)")
	}

	// C. Decrypt the Payload using the recovered AES Session Key
	block, err := aes.NewCipher(sessionKey)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintextBytes, err := gcm.Open(nil, nonceBytes, encDataBytes, nil)
	if err != nil {
		return "", fmt.Errorf("AES decryption failed (Tampered data?)")
	}

	return string(plaintextBytes), nil
}

// --- 3. REPLAY VALIDATION ---
func ValidateRequest(payload string) {
	parts := strings.Split(payload, "|")
	if len(parts) != 4 {
		log.Fatal("Invalid payload structure")
	}

	tsStr, method, secret, endpoint := parts[0], parts[1], parts[2], parts[3]

	// --- MATHEMATICAL REPLAY CHECK ---
	// Convert timestamp string to int64
	reqTime, _ := strconv.ParseInt(tsStr, 10, 64)
	now := time.Now().Unix()

	// Define the "Replay Window" (e.g., 30 seconds)
	// If the request is older than 30s, or "from the future", reject it.
	if now-reqTime > 30 || reqTime > now+5 {
		fmt.Println("❌ REJECTED: Token Expired (Replay Attack Attempt)")
		return
	}

	fmt.Println("✅ ACCEPTED: Token is fresh.")
	fmt.Printf("   Details: %s %s (Secret: %s)\n", method, endpoint, secret)
}

func main() {
	// Setup Keys (Simulated)
	privKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	pubKey := &privKey.PublicKey

	

	// --- SCENARIO ---
	// Construct a payload that might be very long
	longSecret := "ThisSecretIsSuperLongAndWouldNormallyBreakRSAIfWeDidNotUseHybridEncryption_AABBCCDDEEFF"
	payload := fmt.Sprintf("%d|POST|%s|/api/v1/very/secure/resource", time.Now().Unix(), longSecret)

	fmt.Println("1. Encrypting Hybrid Packet...")
	secureToken, err := HybridEncrypt(pubKey, payload)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("   Token Length: %d chars\n", len(secureToken))

	// ... Network Transmission ...

	fmt.Println("2. Decrypting on Server...")
	decryptedPayload, err := HybridDecrypt(privKey, secureToken)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("3. Validating Logic...")
	ValidateRequest(decryptedPayload)

	// --- TEST REPLAY ATTACK ---
	fmt.Println("\n[Simulating a Replay Attack with an old timestamp]")
	oldPayload := fmt.Sprintf("%d|POST|secret|/api", time.Now().Unix()-100) // 100 seconds ago
	ValidateRequest(oldPayload)
}
