package sdk

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)


type RemoteConfig struct {
	PublicKey *rsa.PublicKey
}

type RemotePublicKey struct {
	Pem string `json:"pem"`
}


func (r *RemoteConfig) FetchPublicKey(url string) (error) {
	resp, err := http.Get(url)
	if err != nil {
		return  err
	}
	body := resp.Body
	bytes , err := io.ReadAll(body)
	if err != nil {
		return  err
	}
	var pem *RemotePublicKey
	if err := json.Unmarshal(bytes, &pem); err != nil{
		return  err
	}

	if pem == nil {
		return fmt.Errorf("Error")
	}

	r.PublicKey, err = pem.parseRSAPublicKeyFromPEM()

	if err != nil {
		return err
	}	
	
	return nil

}

func (r *RemotePublicKey) parseRSAPublicKeyFromPEM() (*rsa.PublicKey, error) {
    return parseRSAPublicKeyFromPEM(r.Pem)
}

// Url is the domain or IP of the storage system
// Key is the AES Key used for the symmetric encryption system
type StorageSystemConfig struct {
	Url      string
}

/*
	a. init the Remote Config to fetch the public_key
	b. Make requests
*/
type StorageSystemClient struct {
	Config StorageSystemConfig
	remoteConfig RemoteConfig
}

func (s *StorageSystemClient) Init() error{
	s.remoteConfig = RemoteConfig{}
	if err := s.remoteConfig.FetchPublicKey(s.Config.Url); err != nil {
		return err
	}
	return nil
}






func (s *StorageSystemClient) FetchFiles(jwt *string) (any, error) {
    url := s.Config.Url
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("error creating request: %w", err)
    }

    // 1. Add headers (including the security headers X-Signature, X-Nonce, X-Envelope)
    err = addRequestHeaders(req, jwt, s.remoteConfig.PublicKey)
    if err != nil {
        return nil, fmt.Errorf("error adding request headers: %w", err)
    }

    // 2. Execute the request using the client
    client := &http.Client{Timeout: 30 * time.Second} // Define a client with a timeout
    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("error executing request: %w", err)
    }
    defer resp.Body.Close()

    // 3. Handle the response status
    if resp.StatusCode != http.StatusOK {
        // Read and include the body for better error context
        bodyBytes, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("received non-200 status code: %d, response: %s", resp.StatusCode, string(bodyBytes))
    }

    // 4. Process the successful response body
    var fileList interface{} 
    if err := json.NewDecoder(resp.Body).Decode(&fileList); err != nil {
        return nil, fmt.Errorf("error decoding response body: %w", err)
    }

    // Return the successfully decoded data
    return fileList, nil
}

