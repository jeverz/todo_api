package hsjwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

const headerHS256 string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"

// Validate the JWT and decode the data into a struct
func Decode(data string, o any, secret []byte) error {
	start := strings.Index(data, ".")
	end := strings.LastIndex(data, ".")
	if data[:start] != headerHS256 {
		return fmt.Errorf("invalid header")
	}
	if len(data)-end != 44 {
		return fmt.Errorf("invalid signature length: %v", len(data)-end)
	}
	if hash([]byte(data[:end]), secret) != data[end+1:] {
		return fmt.Errorf("invalid signature")
	}
	payload, err := base64.RawURLEncoding.DecodeString(data[start+1 : end])
	if err != nil {
		return err
	}
	if err := json.Unmarshal(payload, o); err != nil {
		return err
	}
	return nil
}

// Encode the data of a struct into a signed JWT
func Encode(o any, secret []byte) (string, error) {
	header, err := json.Marshal(o)
	if err != nil {
		return "", err
	}
	payload := headerHS256 + "." + base64.RawURLEncoding.EncodeToString(header)
	payload += "." + hash([]byte(payload), secret)
	return payload, nil
}

// Get a raw base64 encoded hash of the data
func hash(data []byte, secret []byte) string {
	sig := hmac.New(sha256.New, secret)
	sig.Write(data)
	return base64.RawURLEncoding.EncodeToString(sig.Sum(nil))
}
