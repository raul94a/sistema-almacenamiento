package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"
)

// SecureEnvelope holds the two parts of our hybrid encryption
type SecureEnvelope struct {
	EncryptedKey  string `json:"k"`  // RSA-encrypted AES key (Base64)
	EncryptedData string `json:"d"`  // AES-encrypted Payload (Base64)
	Nonce         string `json:"n"`  // AES Nonce/IV (Base64)
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
	if now - reqTime > 30 || reqTime > now + 5 {
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