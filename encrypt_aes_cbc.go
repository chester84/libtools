package libtools

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/beego/beego/v2/core/logs"
)

const (
	AesCBCIV string = "7816a1762a9a145c24726ea5790fe39a"
)

const (
	AesPKCS5 int = iota
	AesPKCS7
)

func AesEncryptCBCV2(src string, key string, iv string) (string, error) {
	aesKey := []byte(key[0:16])
	aesIV, _ := hex.DecodeString(iv)
	//logs.Notice("[AesEncryptCBCV2] aesKey: %s, aesIV: %v", aesKey, aesIV)
	ciphertextByte, err := Encrypter([]byte(src), aesKey, aesIV, AesPKCS5)
	if err != nil {
		logs.Warning("[AesEncryptCBCV2] has wrong, err: %v", err)
		return "", err
	}

	return string(Base64Encode(ciphertextByte)), nil
}

func AesDecryptCBCV2(ciphertext string, key string, iv string) (string, error) {
	//logs.Debug("ciphertext:", ciphertext)
	base64Data, err := Base64Decode(ciphertext)
	if err != nil {
		logs.Warning("[AesDecryptCBCV2] base64decode has wrong, err: %v", err)
		return "", err
	}

	aesKey := []byte(key[0:16])
	aesIV, _ := hex.DecodeString(iv)

	ciphertextByte, err := Decrypter(base64Data, aesKey, aesIV, AesPKCS5)
	if err != nil {
		logs.Warning("[AesDecryptCBCV2] call Decrypter has wrong.")
		return "", err
	}
	//logs.Debug("ciphertext:", string(ciphertextByte))

	return string(ciphertextByte), nil
}

func AesDecryptUrlCode(ciphertext string, key string, iv string) (string, error) {
	decodeData, err := UrlDecode(ciphertext)
	if err != nil {
		logs.Warning("urldecode has wrong.")
		return "", err
	}

	return AesDecryptCBC(decodeData, key, iv)
}

func AesDecryptCBC(ciphertext string, key string, iv string) (string, error) {
	base64Data, err := Base64Decode(ciphertext)
	if err != nil {
		logs.Warning("base64decode has wrong, ciphertext: %s", ciphertext)
		return "", err
	}

	if len(key) != 16 ||
		len(iv) != 16 {
		return "", fmt.Errorf("key and iv must be length:16")
	}
	aesKey := []byte(key)
	aesIV := []byte(iv)

	ciphertextByte, err := Decrypter(base64Data, aesKey, aesIV, AesPKCS5)
	if err != nil {
		logs.Warning("call Decrypter has wrong.")
		return "", err
	}
	//logs.Debug("ciphertext:", string(ciphertextByte))

	return string(ciphertextByte), nil
}

//解密
func Decrypter(crypted []byte, key []byte, iv []byte, paddingType int) ([]byte, error) {
	var err error
	emptyBytes := []byte{}

	sourceBlock, err := aes.NewCipher(key)
	if err != nil {
		return emptyBytes, err
	}
	if len(crypted)%sourceBlock.BlockSize() != 0 {
		err = errors.New("crypto/cipher: input not full blocks")
		return emptyBytes, err
	}

	source := make([]byte, len(crypted))
	sourceAes := cipher.NewCBCDecrypter(sourceBlock, iv)
	sourceAes.CryptBlocks(source, crypted)
	if paddingType == AesPKCS5 {
		source = PKCS5UnPadding(source)
	} else {
		source = PKCS7UnPadding(source)
	}

	return source, err
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	check := length - unpadding
	if check < 0 {
		return []byte{}
	}
	return origData[:check]
}

func PKCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	unPadding := int(plantText[length-1])
	return plantText[:(length - unPadding)]
}

// 包装过的简单加密函数
func AesEncryptCBC(src string, key string, iv string) (string, error) {
	//aesKey, _ := hex.DecodeString(key)
	//aesIV, _ := hex.DecodeString(iv)
	if len(key) != 16 ||
		len(iv) != 16 {
		return "", fmt.Errorf("key and iv must be length:16")
	}
	aesKey := []byte(key)
	aesIV := []byte(iv)

	ciphertextByte, err := Encrypter([]byte(src), aesKey, aesIV, AesPKCS5)
	if err != nil {
		return "", err
	}

	return Base64Encode(ciphertextByte), nil
}

//加密
func Encrypter(source []byte, key []byte, iv []byte, paddingType int) ([]byte, error) {
	sourceBlock, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, err
	}
	if paddingType == AesPKCS5 {
		source = PKCS5Padding(source, sourceBlock.BlockSize()) //补全位数，长度必须是 16 的倍数
	} else {
		source = PKCS7Padding(source, sourceBlock.BlockSize())
	}

	sourceCrypted := make([]byte, len(source))
	sourceAes := cipher.NewCBCEncrypter(sourceBlock, iv)
	sourceAes.CryptBlocks(sourceCrypted, source)
	return sourceCrypted, err
}

// 补位
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//url base64编码
func URLBase64Encode(data []byte) string {
	base64encodeBytes := base64.URLEncoding.EncodeToString(data)
	//logs.Debug("base64encode:", base64encodeBytes)
	base64encodeBytes = strings.Replace(base64encodeBytes, "=", ".", -1)
	return base64encodeBytes
}

//url base64解码
func URLBase64Decode(data string) ([]byte, error) {
	data = strings.Replace(data, ".", "=", -1)
	decodeBytes, err := base64.URLEncoding.DecodeString(data)

	return decodeBytes, err
}

//base64编码
func Base64Encode(data []byte) string {
	base64encodeBytes := base64.StdEncoding.EncodeToString(data)
	//logs.Debug("base64encode:", base64encodeBytes)
	return base64encodeBytes
}

//base64解码
func Base64Decode(data string) ([]byte, error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(data)
	return decodeBytes, err
}

//url 编码
func UrlEncode(data string) string {
	encode := url.QueryEscape(data)
	return encode
}

//url 解码
func UrlDecode(data string) (string, error) {
	decodeurl, err := url.QueryUnescape(data)
	return decodeurl, err
}
