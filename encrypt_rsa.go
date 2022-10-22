package libtools

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"

	"github.com/beego/beego/v2/core/logs"
)

var privateKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`)

var publicKey = []byte(`-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAsxxIyT2PftaUkIIHcjuG
cxkypR7iOQKuW6whvMRwp6dQ/FYPLuRaPxGOsYhoSGSs8JnxuO28wyALyQXxL/Rf
q2gR8qq7ADo9W0Dmsg1l6zLqDZtItdMwlcWU69XaUjQgkQsJpLE+xTJBh15FcBi7
mWxd1ykshJAK3f08y2ZygCHAvUtwanNmEPPDPybuZsd7LShduVTZekStreiQV3/u
SNKaT2GaHCCKHlQGHCmMCPNKfzn2gapTKyHbnnAyQ7vrB/lhnxl6gIRQ+0wf6Q4m
R11HpjUVihCvhFvGenn1KTYnYxh/d3fnsU+J4pnG1hD+LpLcSJLIunv2tyRfLPM7
9DTxQjLU/1thVFGdtXPpW7RFwIayPpsoMwDMxW9Lue6SYBVLZQOCsK7RXGH+6wCY
bE/Yip0/BlTlRrlWvIwUGzI/zlPXyMEojM0aA1vxY1pvBZzMGzPubw14SQ2bt0Ox
uXYnybIpKKyct+9CM1UrNvOo2B3X1Avsagqn1U4OcST4hbOZySZtsj/4JX5vlPpA
6Hv1e0BxnXYO/q8cRILyZvDW0rjitPhbDz3TGu9RJSmmpuM9E626mQl7yC6qOmVZ
zg00/5aaD6jazUBZScXF1OYMLo3CdDCyddDEI/biYuZnAfpPyJCiqD28pm6vctcx
sSxmd59oK8f3i1v/7DPkVb8CAwEAAQ==
-----END PUBLIC KEY-----`)

// 加密
func RsaEncrypt(origData []byte) (string, error) {
	//解密pem格式的公钥
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return "", errors.New("public key error")
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	//加密
	bytes, err := rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
	return Base64Encode(bytes), err
}

// 解密
func RsaDecrypt(ciphertext string) (string, error) {
	base64Data, err := Base64Decode(ciphertext)
	if err != nil {
		logs.Warning("[RsaDecrypt] base64decode has wrong, ciphertext: %s", ciphertext)
		return "", err
	}
	//解密
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return "", errors.New("private key error!")
	}
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	// 解密
	btyes, err := rsa.DecryptPKCS1v15(rand.Reader, priv, base64Data)
	if err != nil {
		logs.Warning("call RsaDecrypt has wrong.")
		return "", err
	}
	return string(btyes), nil
}
