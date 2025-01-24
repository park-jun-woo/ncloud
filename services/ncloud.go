package services

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Access struct {
	AccessKey string
	SecretKey string
}

func Request(access *Access, method string, endpoint string, url string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader

	// Convert the body to JSON if it's not nil
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("Failed to marshal request body: %v", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(method, endpoint+url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("Fail to Create HTTP Request: %v", err)
	}

	timestamp := time.Now().UnixMilli()

	signature := makeSignature(access.AccessKey, access.SecretKey, method, url, timestamp)

	req.Header.Set("x-ncp-apigw-timestamp", fmt.Sprintf("%d", timestamp))
	req.Header.Set("x-ncp-iam-access-key", access.AccessKey)
	req.Header.Set("x-ncp-apigw-signature-v2", signature)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Fail to HTTP Request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Fail(CODE=%d): %s %v %s", resp.StatusCode, url, req.Header, string(bodyBytes))
	}

	return resp, nil
}

func makeSignature(accessKeyID string, secretKey string, method string, path string, epochTime int64) string {
	const space = " "
	const newLine = "\n"
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(method))
	h.Write([]byte(space))
	h.Write([]byte(path))
	h.Write([]byte(newLine))
	h.Write([]byte(fmt.Sprintf("%d", epochTime)))
	h.Write([]byte(newLine))
	h.Write([]byte(accessKeyID))
	rawSignature := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(rawSignature)
}
