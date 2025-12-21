package repository

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/google/uuid"
	"storage.client.com/cryptoutils"
	"storage.client.com/header"

	"storage.client.com/repository/models"
)

// El secret es un secreto compartido entre SDK y el sistema de almacenamiento
// Funciona como una capa extra de seguridad, que se adiciona al JWT, Auth,
// timestamp anti replay attacks y el hash del nonce
type StorageRepository struct {
	BaseURL     string
	HTTPClient  *http.Client
	PublicKey   *rsa.PublicKey
	Secret      *string
	cryptoutils cryptoutils.CryptoUtils
}

// Constructor para inicializar el SDK
func NewRepository(baseURL string) *StorageRepository {
	return &StorageRepository{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: time.Second * 30, // Buen hábito en SDKs
		},
		cryptoutils: cryptoutils.CryptoUtils{},
	}
}

func NewRepositoryWithPublicKey(baseUrl string, pk *rsa.PublicKey) *StorageRepository {
	return &StorageRepository{
		BaseURL:   baseUrl,
		PublicKey: pk,
		HTTPClient: &http.Client{
			Timeout: time.Second * 30, // Buen hábito en SDKs
		},
		cryptoutils: cryptoutils.CryptoUtils{},
	}
}

func (c *StorageRepository) addHeaders(jwt *string, req *http.Request) (*http.Request, error) {
	builder := header.RequestHeaderBuilder{}
	aesKey, err := c.cryptoutils.CreateAesKey()
	if err != nil {
		return nil, err
	}
	payload := fmt.Sprintf("%d|%s|%s", time.Now().Unix(), req.Method, req.URL.Path)
	if c.Secret != nil {
		payload = fmt.Sprintf("%s|%s", payload, *c.Secret)
	}
	packet, err := c.cryptoutils.EncryptPayload(aesKey, payload)
	if err != nil {
		return nil, err
	}
	rsaEncryptedAes, err := c.cryptoutils.RsaEncryptAesKey(c.PublicKey,aesKey)
	if err != nil {
		return nil, err
	}
	hexEncryptedAes := c.cryptoutils.HexEncodeEncryptedAesKey(rsaEncryptedAes)
	hexEncodedAesPacket := c.cryptoutils.HexEncodeAesPacket(packet)
	hashNonce := c.cryptoutils.HashNonce(packet.Data)
	builder = *builder.Builder(req).
		AddSite().
		AddDigitalEnvelope(hexEncodedAesPacket.Data).
		AddNonce(hexEncodedAesPacket.Nonce).
		AddAuth(hexEncryptedAes).
		AddHashHeader(hashNonce)

	if jwt != nil {
		builder = *builder.AddJwtHeader(*jwt)
	}

	req = builder.Build()
	return req, nil
}

func (c *StorageRepository) FetchFiles(jwt *string) ([]models.Object, error) {
	// 1. Crear Request
	url := fmt.Sprintf("%s/objects", c.BaseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req, err = c.addHeaders(jwt, req)
	if err != nil {
		return nil, err
	}

	// 3. Ejecutar
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 4. Validar Status Code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error API: status %d", resp.StatusCode)
	}

	// 5. Decodificar JSON a Structs
	var objects []models.Object
	if err := json.NewDecoder(resp.Body).Decode(&objects); err != nil {
		return nil, err
	}

	return objects, nil
}

// FetchFile: Descargar un archivo específico
func (c *StorageRepository) FetchFile(o models.Object, jwt *string) ([]byte, error) {

	url := fmt.Sprintf("%s/objects/%s/download", c.BaseURL, o.ID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req, err = c.addHeaders(jwt, req)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error descargando: status %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// PutObject: Subir archivo nuevo (Multipart)
func (c *StorageRepository) PutObject(creator models.UploadObject, DATA []byte, jwt *string) (uuid.UUID, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	metaDataBytes, _ := json.Marshal(creator)
	_ = writer.WriteField("metadata", string(metaDataBytes))

	part, err := writer.CreateFormFile("file", creator.Filename)
	if err != nil {
		return uuid.Nil, err
	}
	_, err = part.Write(DATA)
	if err != nil {
		return uuid.Nil, err
	}
	writer.Close()

	url := fmt.Sprintf("%s/api/v1/PutObject", c.BaseURL)
	fmt.Printf("URL %s\n", url)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return uuid.Nil, err
	}

	// Aplicar seguridad sobre la request ya formada (incluyendo headers multipart)
	req, err = c.addHeaders(jwt, req)
	if err != nil {
		return uuid.Nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return uuid.Nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return uuid.Nil, fmt.Errorf("error subiendo: status %d", resp.StatusCode)
	}

	var result struct {
		ID uuid.UUID `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return uuid.Nil, err
	}

	return result.ID, nil
}

// RemoveObject: Eliminar archivo
func (c *StorageRepository) RemoveObject(o models.Object, jwt *string) error {
	url := fmt.Sprintf("%s/objects/%s", c.BaseURL, o.ID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req, err = c.addHeaders(jwt, req)
	if err != nil {
		return err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("error eliminando: status %d", resp.StatusCode)
	}

	return nil
}
