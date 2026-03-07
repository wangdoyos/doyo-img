package util

import (
	"crypto/rand"
	"math/big"
)

// GenerateID 使用密码学安全的随机数生成指定长度的 ID
func GenerateID(length int, alphabet string) (string, error) {
	result := make([]byte, length)
	alphabetLen := big.NewInt(int64(len(alphabet)))
	for i := 0; i < length; i++ {
		idx, err := rand.Int(rand.Reader, alphabetLen)
		if err != nil {
			return "", err
		}
		result[i] = alphabet[idx.Int64()]
	}
	return string(result), nil
}

// GenerateDeleteToken 生成带 "tok_" 前缀的 24 位随机删除令牌
func GenerateDeleteToken() (string, error) {
	token, err := GenerateID(24, "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	if err != nil {
		return "", err
	}
	return "tok_" + token, nil
}
