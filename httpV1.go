package libtools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// ContentType 类型定义
type ContentType string

const (
	HttpApplicationJSON        ContentType = "application/json"
	HttpMultipartForm          ContentType = "multipart/form-data"
	HttpApplicationFormEncoded ContentType = "application/x-www-form-urlencoded"
)

// HttpRequest 封装的 HTTP 请求函数，带默认超时 60 秒，允许覆盖超时参数
func HttpRequest(method, urlStr string, headers map[string]string, contentType ContentType, body interface{}, timeout ...time.Duration) ([]byte, int, error) {
	var requestBody io.Reader
	var contentTypeHeader string
	var httpStatusCode int
	var emptyBody []byte

	// 如果用户没有传入超时参数，设置默认超时时间为 60 秒
	var clientTimeout time.Duration
	if len(timeout) > 0 {
		clientTimeout = timeout[0] // 使用传入的超时时间
	} else {
		clientTimeout = 60 * time.Second // 默认 60 秒超时
	}

	switch contentType {
	case HttpApplicationJSON:
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, httpStatusCode, fmt.Errorf("could not marshal json: %v", err)
		}
		requestBody = bytes.NewBuffer(jsonBody)
		contentTypeHeader = string(HttpApplicationJSON)

	case HttpMultipartForm:
		var buffer bytes.Buffer
		writer := multipart.NewWriter(&buffer)

		data := body.(map[string]interface{})
		for key, val := range data {
			switch v := val.(type) {
			case string:
				_ = writer.WriteField(key, v)
			case *os.File:
				file, err := os.Open(v.Name())
				if err != nil {
					return nil, httpStatusCode, fmt.Errorf("could not open file: %v", err)
				}
				defer file.Close()

				part, err := writer.CreateFormFile(key, v.Name())
				if err != nil {
					return nil, httpStatusCode, fmt.Errorf("could not create form file: %v", err)
				}
				_, err = io.Copy(part, file)
				if err != nil {
					return nil, httpStatusCode, fmt.Errorf("could not copy file content: %v", err)
				}
			default:
				return nil, httpStatusCode, fmt.Errorf("unsupported field type: %v", v)
			}
		}

		err := writer.Close()
		if err != nil {
			return nil, httpStatusCode, fmt.Errorf("could not close writer: %v", err)
		}

		requestBody = &buffer
		contentTypeHeader = writer.FormDataContentType()

	case HttpApplicationFormEncoded:
		formData := url.Values{}
		data := body.(map[string]string)
		for key, val := range data {
			formData.Set(key, val)
		}
		requestBody = strings.NewReader(formData.Encode())
		contentTypeHeader = string(HttpApplicationFormEncoded)

	default:
		return nil, httpStatusCode, fmt.Errorf("unsupported content type: %v", contentType)
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest(method, urlStr, requestBody)
	if err != nil {
		return nil, httpStatusCode, fmt.Errorf("could not create http request: %v", err)
	}

	// 设置 Content-Type
	req.Header.Set("Content-Type", contentTypeHeader)

	// 设置自定义的 headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 创建 HTTP 客户端，并设置超时时间
	client := &http.Client{
		Timeout: clientTimeout, // 使用默认或用户提供的超时时间
	}

	// 发送 HTTP 请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, httpStatusCode, fmt.Errorf("could not send http request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return emptyBody, httpStatusCode, err
	}

	return respBody, resp.StatusCode, err
}

// 用法如下
func test() {
	// JSON 请求示例
	jsonHeaders := map[string]string{
		"Authorization": "Bearer YOUR_TOKEN",
		"Custom-Header": "CustomValue",
	}
	jsonBody := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	jsonResponse, _, err := HttpRequest("POST", "https://example.com/json-endpoint", jsonHeaders, HttpApplicationJSON, jsonBody, 10*time.Second)
	fmt.Printf("JSON Response: %s\n", jsonResponse)

	// Multipart/form-data 请求示例
	multipartHeaders := map[string]string{
		"Authorization": "Bearer YOUR_TOKEN",
		"Custom-Header": "CustomValue",
	}
	multipartBody := map[string]interface{}{
		"field1": "value1",
		//"file":   &os.File{Name: "path/to/your/file"},
	}
	multipartResponse, _, err := HttpRequest("POST", "https://example.com/upload", multipartHeaders, HttpMultipartForm, multipartBody, 10*time.Second)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Multipart Response: %s\n", multipartResponse)

	// Application/x-www-form-urlencoded 请求示例
	formHeaders := map[string]string{
		"Authorization": "Bearer YOUR_TOKEN",
		"Custom-Header": "CustomValue",
	}
	formBody := map[string]string{
		"field1": "value1",
		"field2": "value2",
	}
	formResponse, _, err := HttpRequest("POST", "https://example.com/form-endpoint", formHeaders, HttpApplicationFormEncoded, formBody)
	fmt.Printf("Form Response: %s\n", formResponse)
}
