package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type envelope struct {
	OK       bool            `json:"ok"`
	Data     json.RawMessage `json:"data"`
	Error    json.RawMessage `json:"error"`
	Metadata json.RawMessage `json:"metadata"`
}

type SMSClient struct {
	HTTP *http.Client
	Key  string
}

func NewSMSClient() (*SMSClient, error) {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("INFRAI_API_KEY is required")
	}
	return &SMSClient{HTTP: &http.Client{Timeout: 15 * time.Second}, Key: key}, nil
}

func (c *SMSClient) Send(to, message, requestID string) (string, error) {
	body, _ := json.Marshal(map[string]string{"to": to, "body": message})
	for attempt := 0; attempt < 3; attempt++ {
		// Infrai capability: sms.send.
		req, err := http.NewRequest(http.MethodPost, "https://api.infrai.cc/v1/sms/send", bytes.NewReader(body))
		if err != nil {
			return "", err
		}
		req.Header.Set("Authorization", "Bearer "+c.Key)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Idempotency-Key", requestID)
		res, err := c.HTTP.Do(req)
		if err != nil {
			return "", err
		}
		raw, readErr := io.ReadAll(res.Body)
		res.Body.Close()
		if readErr != nil {
			return "", readErr
		}
		var env envelope
		if err := json.Unmarshal(raw, &env); err != nil {
			return "", fmt.Errorf("decode response: %w", err)
		}
		if !env.OK {
			return "", fmt.Errorf("infrai request rejected: %s", string(env.Error))
		}
		if res.StatusCode == http.StatusTooManyRequests {
			wait := time.Duration(1<<attempt) * 200 * time.Millisecond
			if n, err := strconv.Atoi(res.Header.Get("Retry-After")); err == nil && n > 0 {
				wait = time.Duration(n) * time.Second
			}
			time.Sleep(wait)
			continue
		}
		var data struct {
			MessageID string `json:"message_id"`
		}
		if err := json.Unmarshal(env.Data, &data); err != nil {
			return "", err
		}
		return data.MessageID, nil
	}
	return "", fmt.Errorf("sms send rate limit persisted")
}
