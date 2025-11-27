package libtools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
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

	// 设置默认 60 秒超时
	clientTimeout := 60 * time.Second
	if len(timeout) > 0 {
		clientTimeout = timeout[0]
	}

	switch contentType {
	case HttpApplicationJSON:
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, httpStatusCode, fmt.Errorf("could not marshal json: %w", err)
		}
		requestBody = bytes.NewBuffer(jsonBody)
		contentTypeHeader = string(HttpApplicationJSON)

	case HttpMultipartForm:
		// body 必须是 map[string]interface{}
		data, ok := body.(map[string]interface{})
		if !ok {
			return nil, httpStatusCode, fmt.Errorf("HttpMultipartForm expects body of type map[string]interface{}")
		}

		var buffer bytes.Buffer
		writer := multipart.NewWriter(&buffer)

		for key, val := range data {
			switch v := val.(type) {
			case string:
				if err := writer.WriteField(key, v); err != nil {
					return nil, httpStatusCode, fmt.Errorf("could not write field %s: %w", key, err)
				}

			case *os.File:
				// 不在这里 Close()，由调用方负责关闭文件句柄
				if _, err := v.Seek(0, io.SeekStart); err != nil {
					// 非致命，但尽量回到文件头
					// 如果 Seek 失败，仍然尝试读取
				}
				part, err := writer.CreateFormFile(key, filepath.Base(v.Name()))
				if err != nil {
					return nil, httpStatusCode, fmt.Errorf("could not create form file for %s: %w", key, err)
				}
				if _, err := io.Copy(part, v); err != nil {
					return nil, httpStatusCode, fmt.Errorf("could not copy file content for %s: %w", key, err)
				}

			case io.Reader:
				// 支持任意 io.Reader（例如 bytes.Buffer、bytes.Reader）
				part, err := writer.CreateFormFile(key, key) // 如果没有文件名，用 key 作为占位名
				if err != nil {
					return nil, httpStatusCode, fmt.Errorf("could not create form file for reader %s: %w", key, err)
				}
				if _, err := io.Copy(part, v); err != nil {
					return nil, httpStatusCode, fmt.Errorf("could not copy reader content for %s: %w", key, err)
				}

			default:
				return nil, httpStatusCode, fmt.Errorf("unsupported field type for key %s: %T", key, v)
			}
		}

		if err := writer.Close(); err != nil {
			return nil, httpStatusCode, fmt.Errorf("could not close multipart writer: %w", err)
		}

		requestBody = &buffer
		contentTypeHeader = writer.FormDataContentType()

	case HttpApplicationFormEncoded:
		formData := url.Values{}
		data, ok := body.(map[string]string)
		if !ok {
			return nil, httpStatusCode, fmt.Errorf("HttpApplicationFormEncoded expects body of type map[string]string")
		}
		for key, val := range data {
			formData.Set(key, val)
		}
		requestBody = strings.NewReader(formData.Encode())
		contentTypeHeader = string(HttpApplicationFormEncoded)

	case HttpRawBody:
		// 支持 []byte, *bytes.Buffer, io.Reader
		switch v := body.(type) {
		case []byte:
			requestBody = bytes.NewReader(v)
		case *bytes.Buffer:
			requestBody = v
		case io.Reader:
			requestBody = v
		default:
			return nil, httpStatusCode, fmt.Errorf("HttpRawBody only accepts []byte, *bytes.Buffer or io.Reader, got %T", body)
		}
		// contentTypeHeader 留空，由调用者在 headers 中手动设置

	default:
		return nil, httpStatusCode, fmt.Errorf("unsupported content type: %v", contentType)
	}

	// 构建 request
	req, err := http.NewRequest(method, urlStr, requestBody)
	if err != nil {
		return nil, httpStatusCode, fmt.Errorf("could not create http request: %w", err)
	}

	// 只有在非 RawBody 情况下，才自动设置 Content-Type
	if contentTypeHeader != "" {
		req.Header.Set("Content-Type", contentTypeHeader)
	}

	// 用户 Header 覆盖（包含可能的 Content-Type）
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// 使用 context 以便可扩展取消（可选）
	ctx, cancel := context.WithTimeout(req.Context(), clientTimeout)
	defer cancel()
	req = req.WithContext(ctx)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, httpStatusCode, fmt.Errorf("could not send http request: %w", err)
	}
	defer resp.Body.Close()

	httpStatusCode = resp.StatusCode

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, httpStatusCode, fmt.Errorf("could not read response body: %w", err)
	}

	return respBody, httpStatusCode, nil
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
