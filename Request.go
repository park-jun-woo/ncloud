package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

func Request(access Access, method string, uri string) (*http.Response, error) {
	req, err := http.NewRequest(method, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("Fail to Create HTTP Request: %v", err)
	}

	timestamp := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)

	message := fmt.Sprintf("%s %s\n%s\n%s", method, uri, timestamp, access.AccessKey)
	mac := hmac.New(sha256.New, []byte(access.SecretKey))
	mac.Write([]byte(message))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	req.Header.Set("x-ncp-apigw-timestamp", timestamp)
	req.Header.Set("x-ncp-iam-access-key", access.AccessKey)
	req.Header.Set("x-ncp-apigw-signature-v2", signature)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Fail to HTTP Request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Fail(CODE=%d): %s %s %s", resp.StatusCode, method, uri, string(bodyBytes))
	}

	return resp, nil
}
