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
	HttpRawBody                ContentType = "raw" // 新增，用于手动构造 body
)

// HttpRequest 封装的 HTTP 请求函数，带默认超时 60 秒，允许覆盖超时参数
func HttpRequest(method, urlStr string, headers map[string]string, contentType ContentType, body interface{}, timeout ...time.Duration) ([]byte, int, error) {
	var requestBody io.Reader
	var contentTypeHeader string
	var httpStatusCode int
	var emptyBody []byte

	// 设置默认 60 秒超时
	clientTimeout := 60 * time.Second
	if len(timeout) > 0 {
		clientTimeout = timeout[0]
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
		// ⚠️ 仅适合 map 自动构建 multipart 的情况
		var buffer bytes.Buffer
		writer := multipart.NewWriter(&buffer)

		data := body.(map[string]interface{})
		for key, val := range data {
			switch v := val.(type) {
			case string:
				_ = writer.WriteField(key, v)

			case *os.File:
				part, err := writer.CreateFormFile(key, filepath.Base(v.Name()))
				if err != nil {
					return nil, httpStatusCode, fmt.Errorf("could not create form file: %v", err)
				}

				_, err = io.Copy(part, v)
				if err != nil {
					return nil, httpStatusCode, fmt.Errorf("could not copy file content: %v", err)
				}

			default:
				return nil, httpStatusCode, fmt.Errorf("unsupported field type: %v", v)
			}
		}

		_ = writer.Close()
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

	case HttpRawBody:
		// 🚀 这里 body 必须是 []byte 或 bytes.Buffer
		switch v := body.(type) {
		case []byte:
			requestBody = bytes.NewReader(v)
		case *bytes.Buffer:
			requestBody = v
		default:
			return nil, httpStatusCode, fmt.Errorf("HttpRawBody only accepts []byte or *bytes.Buffer")
		}

		// Content-Type 由调用者自行设置，不能自动覆盖
		contentTypeHeader = "" // 标记不自动设置

	default:
		return nil, httpStatusCode, fmt.Errorf("unsupported content type: %v", contentType)
	}

	// ---------------------------
	// 构建 request
	// ---------------------------
	req, err := http.NewRequest(method, urlStr, requestBody)
	if err != nil {
		return nil, httpStatusCode, fmt.Errorf("could not create http request: %v", err)
	}

	// 只有在非 RawBody 情况下，才自动设置 Content-Type
	if contentTypeHeader != "" {
		req.Header.Set("Content-Type", contentTypeHeader)
	}

	// 用户 Header 永远最后覆盖
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: clientTimeout}

	resp, err := client.Do(req)
	if err != nil {
		return nil, httpStatusCode, fmt.Errorf("could not send http request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
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
