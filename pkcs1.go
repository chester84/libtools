package libtools

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"

	"github.com/beego/beego/v2/core/logs"
)

func PKCS1InitPubKey(key []byte) (pub *rsa.PublicKey, err error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return pub, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return pub, err
	}
	pub = pubInterface.(*rsa.PublicKey)
	return
}

func PKCS1Verify(key []byte, data string, sign string) error {

	pub, _ := PKCS1InitPubKey(key)

	h := crypto.MD5.New()
	h.Write([]byte(data))
	hashed := h.Sum(nil)

	//decodedSign, err := base64.RawURLEncoding.DecodeString(sign)
	decodedSign, err := Base64Decode(sign)
	if err != nil {
		return err
	}

	return rsa.VerifyPKCS1v15(pub, crypto.MD5, hashed, decodedSign)
}

func PKCS1RsaEncrypt(data []byte, key []byte) ([]byte, error) {
	pub, _ := PKCS1InitPubKey(key)
	partLen := pub.N.BitLen()/8 - 11
	chunks := split(data, partLen)

	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		bytes, err := rsa.EncryptPKCS1v15(rand.Reader, pub, chunk)
		if err != nil {
			return bytes, err
		}
		buffer.Write(bytes)
	}

	return buffer.Bytes(), nil
}

func PKCS1InitPubKeyPrivateDecrypt(encrypted string, pubkey string, privateKey []byte) (string, error) {
	block, _ := pem.Decode([]byte(pubkey))
	if block == nil {
		return "", errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	pub := pubInterface.(*rsa.PublicKey)
	partLen := pub.N.BitLen() / 8
	//raw, err := base64.RawURLEncoding.DecodeString(encrypted)
	raw, err := Base64Decode(encrypted)
	chunks := split([]byte(raw), partLen)

	block, _ = pem.Decode(privateKey)
	if block == nil {
		return "", errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	pri, ok := priv.(*rsa.PrivateKey)

	logs.Debug("priv", ok)

	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, pri, chunk)
		if err != nil {
			return "", err
		}
		buffer.Write(decrypted)
	}

	return buffer.String(), err
}
