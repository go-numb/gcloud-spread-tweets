package models

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"golang.org/x/crypto/bcrypt"
)

func ToHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// KEYSIZE AES-256-GCMに必要なキーのサイズ
const KEYSIZE = 32

// GenerateKey は、AES-256-GCMに必要なキーを生成します
func GenerateKey() ([]byte, error) {
	key := make([]byte, KEYSIZE)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}
	return key, nil
}

// EncryptPassword は、与えられたパスワードを暗号化します
func EncryptPassword(password string, key []byte) (string, error) {
	// 新しいAES暗号化ブロックを作成
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// GCMモードの暗号化オブジェクトを作成
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// ノンスを生成
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// パスワードを暗号化
	ciphertext := gcm.Seal(nonce, nonce, []byte(password), nil)

	// 暗号文をbase64でエンコード
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptPassword は、与えられた暗号文を復号化します
func DecryptPassword(ciphertext string, key []byte) (string, error) {
	// 暗号文をbase64でデコード
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	// 新しいAES暗号化ブロックを作成
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// GCMモードの暗号化オブジェクトを作成
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// ノンスのサイズを取得
	nonceSize := gcm.NonceSize()
	if len(ciphertextBytes) < nonceSize {
		return "", errors.New("invalid ciphertext")
	}

	// ノンスと暗号文を分離
	nonce, ciphertextBytes := ciphertextBytes[:nonceSize], ciphertextBytes[nonceSize:]

	// 復号化
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
