package sdk

import (
	"crypto/sha3"
	"encoding/base64"
	"encoding/hex"
)

type aesKeySize int

const (
	AES_16_BYTES_KEY aesKeySize = 16
	AES_24_BYTES_KEY aesKeySize = 24
	AES_32_BYTES_KEY aesKeySize = 32
)

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