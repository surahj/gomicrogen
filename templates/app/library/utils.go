package library

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"math"
)

func ToFixed(num float64, precision int) float64 {

	output := math.Pow(10, float64(precision))

	return float64(round(num*output)) / output
}

func Round(num float64, nbDigits float64) float64 {

	pow := math.Pow(10., nbDigits)
	rounded := float64(int(num*pow)) / pow

	return rounded
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

// GenerateSignature calculates an HMAC-SHA256 signature for the given request body using AUTH_TOKEN.
func GetSignature(key string, body string) string {
	return computeHmacSha256(body, key)
}

func computeHmacSha256(message, secret string) string {

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))

	res := hex.EncodeToString(h.Sum(nil))

	log.Printf("HMAC-SHA256 of %s with secret %s = %s", message, secret, res)

	return res
}

func NormalizeJSON(input []byte) (string, error) {
	var raw json.RawMessage // Preserve the raw JSON structure
	err := json.Unmarshal(input, &raw)
	if err != nil {
		return "", err
	}

	// Re-marshal without changing the order
	var buf bytes.Buffer
	err = json.Compact(&buf, raw) // Removes unnecessary spaces and newlines
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

