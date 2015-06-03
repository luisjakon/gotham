package gotham

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"strconv"
	"time"
)

var SecretKey = []byte(NewSecretKey(16))

var (
	Sign   = sign
	Verify = verify
)

func NewSecretKey(length int) string {
	return generateRandomKey(length)
}

func sign(message string) string {
	return hmacSHA256(message)
}

func verify(m1, m2 string) bool {
	return m1 == hmacSHA256(m2)
}

func timestamp(message string) string {
	expiration := strconv.FormatInt(time.Now().Unix(), 10)
	return expiration + "::" + sign(expiration) + "::" + sign(message)
}

func generateRandomKey(length int) string {
	k := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(k)[:length]
}

func hmacSHA256(message string) string {
	h := hmac.New(sha256.New, SecretKey)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
