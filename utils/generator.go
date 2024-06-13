package utils

import (
	"crypto/sha256"
	"encoding/binary"
)

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func GenerateShortLink(originalUrl string) string {
	hash := sha256.Sum256([]byte(originalUrl))
	uniqueId := binary.BigEndian.Uint64(hash[:8])
	return base62Encode(uniqueId)
}

func base62Encode(num uint64) string {
	if num == 0 {
		return string(base62Chars[0])
	}

	result := ""
	for num > 0 {
		remainder := num % 62
		result = string(base62Chars[remainder]) + result
		num = num / 62
	}

	if len(result) > 7 {
		return result[:7]
	}
	for len(result) < 7 {
		// 不足7位的前面补0
		result = "0" + result
	}
	return result
}
