package libtools

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"

	"github.com/beego/beego/v2/core/logs"
)

//md5方法
func Md5(s string) string {
	return Md5Bytes([]byte(s))
}

func Md5Bytes(buf []byte) string {
	h := md5.New()
	h.Write(buf)
	return hex.EncodeToString(h.Sum(nil))
}

func Sha1(data string) string {
	shaFn := sha1.New()
	shaFn.Write([]byte(data))
	return hex.EncodeToString(shaFn.Sum([]byte("")))
}

func HmacSha256(date, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(date))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
	//	hex.EncodeToString(h.Sum(nil))
	// return base64.StdEncoding.EncodeToString([]byte(sha))
}

func Sha256(data string) string {
	shaFn := sha256.New()
	shaFn.Write([]byte(data))
	return hex.EncodeToString(shaFn.Sum([]byte("")))
}

func HmacSHA1(key string, data string) string {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}

//Guid方法
func GetGuid() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return Md5(base64.URLEncoding.EncodeToString(b))
}

func PasswordEncrypt(password string, salt int64) string {
	saltStr := strconv.FormatInt(salt, 10)
	md51 := Md5(password + "$" + saltStr)
	md52 := Md5(saltStr)

	result1 := SubString(md52, 24, 8)
	result2 := SubString(md51, 0, 24)

	return Md5(result1 + result2)
}

func Stringify(obj interface{}) string {
	type stringer interface {
		String() string
	}

	switch obj := obj.(type) {
	case stringer:
		return obj.String()

	case string:
		return obj

	case int:
		return Int2Str(obj)

	case int64:
		return Int642Str(obj)

	case float64:
		return Float642Str(obj)

	case bool:
		if obj {
			return "true"
		}
		return "false"

	default:
		logs.Error("[Stringify] no match, obj: %#v", obj)
		return ""
	}
}

// Signature 空参数值的参数名参与签名
func Signature(params map[string]interface{}, secret string) string {
	paramLen := len(params)
	if paramLen <= 0 {
		logs.Warning("[Signature] params len is 0, params: %v", params)
		return ""
	}

	cntr := make([]string, paramLen)
	var i int = 0
	for k, _ := range params {
		cntr[i] = k
		i++
	}

	// 按字典序列排序
	sort.Strings(cntr)

	str := "" // 待签名字符串
	for i = 0; i < paramLen; i++ {
		key := cntr[i]
		str += fmt.Sprintf("%s=%s&", key, Stringify(params[key]))
	}
	str += secret
	logs.Debug("[Signature] need signature str:", str)

	return Md5(str)
}

func Signature4Struct(obj interface{}, secret string) string {
	bson, err := json.Marshal(obj)
	if err != nil {
		logs.Error("[Signature4Struct] struct2json encdoe has wrong, err: %v", err)
		return ""
	}

	var params map[string]interface{}
	errDc := json.Unmarshal(bson, &params)
	if errDc != nil {
		logs.Error("[Signature4Struct] json2map has wrong, errDc: %v", errDc)
		return ""
	}

	delete(params, "signature")
	delete(params, "sign")

	return Signature(params, secret)
}

func Signature4StructV2(obj interface{}, secret string) string {

	params := Struct2Map(obj)

	// 显示删除签名参数名
	delete(params, "signature")
	delete(params, "sign")

	return SignatureV2(params, secret)
}

// SignatureV2 空参数不参与签名
func SignatureV2(params map[string]interface{}, secret string) string {
	paramLen := len(params)
	if paramLen <= 0 {
		logs.Warning("[SignatureV2] params len is 0, params: %v", params)
		return ""
	}

	cntr := make([]string, paramLen)
	var i int = 0
	for k, _ := range params {
		cntr[i] = k
		i++
	}

	// 按字典序列排序
	sort.Strings(cntr)

	str := "" // 待签名字符串
	for i = 0; i < paramLen; i++ {
		key := cntr[i]
		value := Stringify(params[key])
		if len(value) > 0 {
			str += fmt.Sprintf("%s=%s&", key, value)
		}
	}
	str += secret
	logs.Debug("[SignatureV2] need signature str:", str)

	return Md5(str)
}
