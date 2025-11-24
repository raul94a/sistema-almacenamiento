package cryptoutils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha3"
	"encoding/hex"
	"fmt"
	"io"
	"storage.client.com/aes_packet"
)

type CryptoUtils struct {
}

// ---------------------------------------------------------------------
// Core Cryptographic Operations
// ---------------------------------------------------------------------

// CreateAesKey generates a new symmetric AES key (and potentially the nonce/IV).
func (c *CryptoUtils) CreateAesKey() (key []byte, err error) {
	const AES_KEY_SIZE int = 32
    // Implementation would generate a secure random 256-bit AES key.
    aesKey := make([]byte, AES_KEY_SIZE)
	length, err := io.ReadFull(rand.Reader, aesKey)
	if err != nil {
		// Handle the error, as the buffer was not completely filled.
		// For a cryptographic operation like this, a non-nil error is a critical failure.
		return nil, fmt.Errorf("error reading random bytes: %w", err)
	}
	if AES_KEY_SIZE != length {
		return nil, fmt.Errorf("read length does not match size")
	}
	return aesKey, nil
}

// EncryptPayload encrypts the string payload using the provided AES key.
// It returns an AesPacket containing the ciphertext and the generated nonce.
func (c *CryptoUtils) EncryptPayload(key []byte, payload string) (*aes_packet.AesPacket, error) {
    // Implementation would use AES-GCM or similar mode.
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
	aesPacket := &aes_packet.AesPacket{
		Nonce: nonce,
		Data:  cipherData,
	}
	return aesPacket, nil
}

// RsaEncryptAesKey encrypts the symmetric AES key using the storage server's
// RSA Public Key (creating the "Auth" header value).
func (c *CryptoUtils) RsaEncryptAesKey(pubKey *rsa.PublicKey, aesKey []byte) ([]byte, error) {
	encryptedKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pubKey, aesKey, nil)
	if err != nil {
		return nil, err
	}
	return encryptedKey,nil
}

// ---------------------------------------------------------------------
// Encoding and Hashing Utilities
// ---------------------------------------------------------------------

// HashNonce calculates the SHA3-256 hash of the nonce, hex-encoding the result.
func (c *CryptoUtils) HashNonce(nonce []byte) (string) {
    // Implementation would perform SHA3-256 hash and then hex encoding.
    h := sha3.New256()
	h.Write(nonce)
	digest := h.Sum(nil)
	hexEncoding := hex.EncodeToString(digest)
	return hexEncoding
}

// HexEncodeAesPacket hex-encodes the byte slices within an AesPacket
// for transmission as string headers (X-Digital-Envelope and X-Nonce).
type HexEncodeNonce string
type HexEncodeCiphertext string
func (c *CryptoUtils) HexEncodeAesPacket(p *aes_packet.AesPacket) *aes_packet.HexEncodedAesPacket {
	return p.EncodeToHex()
}

type aesKey []byte
// HexEncodeEncryptedAesKey hex-encodes the RSA-encrypted AES key
// (for transmission as the X-Auth header).
func (c *CryptoUtils) HexEncodeEncryptedAesKey(encryptedKey aesKey) string {
	return hex.EncodeToString(encryptedKey)
}