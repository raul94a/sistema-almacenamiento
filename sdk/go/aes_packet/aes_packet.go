package aes_packet

import "encoding/hex"

type AesPacket struct {
	Data []byte
	Nonce []byte
}

func (a *AesPacket) EncodeToHex() *HexEncodedAesPacket {
	data := hex.EncodeToString(a.Data)
	nonce := hex.EncodeToString(a.Nonce)
	return &HexEncodedAesPacket{
		Data: data,
		Nonce: nonce,
	}
}

type HexEncodedAesPacket struct {
	Data string
	Nonce string
}