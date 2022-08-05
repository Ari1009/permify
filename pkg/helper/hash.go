package helper

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5Hash */
func MD5Hash(text string) string {
	hash := md5.New()
	hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}