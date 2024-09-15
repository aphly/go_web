package crypt

import (
	"crypto/sha256"
	"encoding/hex"
)

func ShaEn(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	hashedPassword := hash.Sum(nil)
	return hex.EncodeToString(hashedPassword)
}
