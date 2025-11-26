package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"
)

// WordList 包含 数字、大小写字母的字符集
var WordList []string

func init() {
	// WORD_LIST
	for ind := 48; ind < 58; ind++ { // 数字
		WordList = append(WordList, string(rune(ind)))
	}
	for ind := 65; ind < 91; ind++ { // 大写字母
		WordList = append(WordList, string(rune(ind)))
	}
	for ind := 97; ind < 123; ind++ { // 小写字母
		WordList = append(WordList, string(rune(ind)))
	}
}

func PKCS5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}

func EncryptCBC(plaintext []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("invalid decrypt key")
	}
	blockSize := block.BlockSize()
	plaintext = PKCS5Padding(plaintext, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintext))
	blockMode.CryptBlocks(ciphertext, plaintext)
	return ciphertext, nil
}

func DecryptCBC(ciphertext []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("invalid decrypt key")
	}
	blockSize := block.BlockSize()
	if len(ciphertext) < blockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	if len(ciphertext)%blockSize != 0 {
		return nil, fmt.Errorf("ciphertext is not a multiple of the block size")
	}
	blockModel := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	blockModel.CryptBlocks(plaintext, ciphertext)
	plaintext = PKCS5UnPadding(plaintext)
	return plaintext, nil
}

func CheckAesKey(strKey string) []byte {
	keyLen := len(strKey)
	arrKey := []byte(strKey)
	if keyLen >= 32 {
		return arrKey[:32]
	}
	if keyLen >= 24 {
		return arrKey[:24]
	}
	if keyLen >= 16 {
		return arrKey[:16]
	}
	tmp := make([]byte, 16)
	for i := 0; i < 16; i++ {
		if i < keyLen {
			tmp[i] = arrKey[i]
		} else {
			tmp[i] = '0'
		}
	}
	return tmp
}

func DecodeAesString(aesStr string) (string, error) {
	if aesStr == "" {
		return "", nil
	}
	strByte, err := base64.StdEncoding.DecodeString(aesStr)
	if err != nil {
		return "", err
	}
	strAes, err := DecryptCBC(strByte, CheckAesKey(SecretKeySetting.AesKey), CheckAesKey(SecretKeySetting.AesVi))
	if err != nil {
		return "", err
	}
	return string(strAes), nil
}

func EncodeAesString(str string) (string, error) {
	aesStr, err := EncryptCBC([]byte(str), CheckAesKey(SecretKeySetting.AesKey), CheckAesKey(SecretKeySetting.AesVi))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(aesStr), nil
}

// HideNickName 隐藏真实nickname
func HideNickName(nickname string) string {
	r := []rune(nickname)
	if len(r) == 0 {
		return "***"
	}
	if len(r) == 1 {
		return string(r[0:1]) + "***"
	}
	if len(r) <= 3 {
		return string(r[0:1]) + "****" + string(r[len(r)-1:])
	}
	return string(r[0:1]) + "****" + string(r[len(r)-1:])
}

// GetMockNickName 生成一个假昵称(含脱敏)
func GetMockNickName() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	first := WordList[r.Intn(len(WordList))]
	last := WordList[r.Intn(len(WordList))]

	return first + "****" + last
}
