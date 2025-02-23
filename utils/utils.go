package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/szwtdl/jxtxedu/types"
	"io"
	"strconv"

	"golang.org/x/crypto/pbkdf2"
)

const (
	password   = "afNpD!r9yW@Ct4jZ"
	iterations = 999
	keySize    = 64 // AES-256
)

// pkcs7Pad 填充数据到 blockSize 长度
func pkcs7Pad(data []byte, blockSize int) []byte {
	// 计算填充长度
	padding := blockSize - len(data)%blockSize
	// 使用重复的字节填充数据
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// pkcs7Unpad 去除 PKCS7 填充
func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	// 如果数据为空或数据长度不是 blockSize 的倍数，返回错误
	if len(data) == 0 || len(data)%blockSize != 0 {
		return nil, errors.New("invalid padding size")
	}

	// 获取最后一个字节作为填充长度
	paddingLen := int(data[len(data)-1])

	// 检查填充长度是否合理
	if paddingLen == 0 || paddingLen > blockSize {
		return nil, errors.New("invalid padding")
	}

	// 确保填充长度与最后几个字节匹配
	for i := len(data) - paddingLen; i < len(data)-1; i++ {
		if data[i] != byte(paddingLen) {
			return nil, errors.New("invalid padding")
		}
	}

	// 返回去除填充后的数据
	return data[:len(data)-paddingLen], nil
}

// 加密数据

func EncryptData(data interface{}) (string, error) {
	var dataStr string
	switch v := data.(type) {
	case string:
		dataStr = v
	default:
		jsonData, err := json.Marshal(data)
		if err != nil {
			return "", err
		}
		dataStr = string(jsonData)
	}

	// 生成 256 位（32字节）的随机 salt
	salt := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}

	// 生成 128 位（16字节）的随机 IV
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// 使用 PBKDF2 生成 512 位（64字节）密钥
	key := pbkdf2.Key([]byte(password), salt, iterations, keySize, sha512.New)
	// 创建 AES 加密块
	block, err := aes.NewCipher(key[:32]) // 取前 32 字节作为 AES 密钥
	if err != nil {
		return "", err
	}

	// 使用 CBC 模式加密
	paddedData := pkcs7Pad([]byte(dataStr), aes.BlockSize)
	ciphertext := make([]byte, len(paddedData))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, paddedData)

	// 将加密数据、salt 和 iv 转换为 JSON
	encryptedData := map[string]string{
		"ciphertext": base64.StdEncoding.EncodeToString(ciphertext),
		"salt":       hex.EncodeToString(salt),
		"iv":         hex.EncodeToString(iv),
	}

	// 转换为 JSON 格式返回
	encryptedJSON, err := json.Marshal(encryptedData)
	if err != nil {
		return "", err
	}

	return string(encryptedJSON), nil
}

// 解密数据

func DecryptData(data map[string]string) (string, error) {
	// 解析 salt
	salt, err := hex.DecodeString(data["salt"])
	if err != nil {
		return "", err
	}

	// 解析 iv
	iv, err := hex.DecodeString(data["iv"])
	if err != nil {
		return "", err
	}

	// 解析 ciphertext
	ciphertext, err := base64.StdEncoding.DecodeString(data["ciphertext"])
	if err != nil {
		return "", err
	}

	// 使用 PBKDF2 生成 512 位（64字节）密钥
	key := pbkdf2.Key([]byte(password), salt, iterations, keySize, sha512.New)

	// 创建 AES 解密块
	block, err := aes.NewCipher(key[:32]) // 取前 32 字节作为 AES 密钥
	if err != nil {
		return "", err
	}

	// 确保密文长度为块大小的整数倍
	if len(ciphertext)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of the block size")
	}

	// 使用 CBC 模式解密
	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	// 去除 PKCS#7 填充
	plaintext, err = pkcs7Unpad(plaintext, aes.BlockSize)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func DecryptEncryptedData(responseApi types.ResponseApi) (string, error) {
	// 获取加密数据
	responseData, ok := responseApi.Data.(string)
	if !ok {
		return "", errors.New("response data is not a string")
	}

	// 解析加密数据
	var encryptedData types.EncryptedData
	if err := json.Unmarshal([]byte(responseData), &encryptedData); err != nil {
		return "", errors.New("解析加密数据: " + err.Error())
	}

	// 验证加密数据是否完整
	if encryptedData.Ciphertext == "" || encryptedData.IV == "" || encryptedData.Salt == "" {
		return "", errors.New("encrypted data is empty")
	}

	// 解密数据
	encryptData := map[string]string{
		"ciphertext": encryptedData.Ciphertext,
		"salt":       encryptedData.Salt,
		"iv":         encryptedData.IV,
	}
	decryptData, err := DecryptData(encryptData)
	if err != nil {
		return "", err
	}

	return decryptData, nil
}

func JsonUnmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func JsonMarshal(v interface{}) []byte {
	data, _ := json.Marshal(v)
	return data
}

func ToInt(value interface{}) (int, error) {
	switch v := value.(type) {
	case string:
		if floatVal, err := strconv.ParseFloat(v, 64); err == nil {
			// 将字符串解析为 float64，再转为 int
			return int(floatVal), nil
		}
		return strconv.Atoi(v)
	case int, int8, int16, int32, int64:
		return v.(int), nil
	case float32, float64:
		return int(v.(float64)), nil
	default:
		return 0, fmt.Errorf("invalid type")
	}
}

func ToString(value interface{}) string {
	return fmt.Sprintf("%v", value)
}
