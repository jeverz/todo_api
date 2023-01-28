package utils

import (
	"crypto/sha256"
	"fmt"
)

func ShaEncode(s *string) *string {
	h := sha256.New()
	h.Write([]byte(*s))
	str := fmt.Sprintf("%x", h.Sum(nil))
	return &str
}
