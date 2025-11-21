package sdk

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
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
