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

	result := make([]byte, 0)
	for num > 0 {
		remainder := num % 62
		result = append([]byte{base62Chars[remainder]}, result...)
		num = num / 62
	}

	// 确保结果长度为 7 位，不足则前面补 0
	for len(result) < 7 {
		result = append([]byte{'0'}, result...)
	}

	// 如果结果长度超过 7 位，则截取前 7 位
	if len(result) > 7 {
		return string(result[:7])
	}

	return string(result)
}
