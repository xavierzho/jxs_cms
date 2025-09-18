package util

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func HttpPost(urlStr string, body map[string]interface{}, header map[string]string) (map[string]interface{}, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("HTTP Post json.Marshal %s %+v error: %s", urlStr, body, err.Error())
	}
	request, err := http.NewRequest("POST", urlStr, strings.NewReader(string(b)))
	if err != nil {
		return nil, fmt.Errorf("HTTP NewRequest POST %s error: %s", urlStr, err.Error())
	}

	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	// 添加额外的请求头
	for key, val := range header {
		request.Header.Set(key, val)
	}
	ctx, cancel := context.WithTimeout(request.Context(), time.Second*60)
	defer cancel()
	request = request.WithContext(ctx)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("HTTP POST %s error: %+v %s", urlStr, request, err.Error())
	}
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("HTTP POST RealAll [%s] error: %+v %s", urlStr, resp.Body, err.Error())
	}

	params := make(map[string]interface{})
	// json转map
	err = json.Unmarshal(respBytes, &params)
	if err != nil {
		return nil, fmt.Errorf("HTTP Unmarshal [%s] error: %s", string(respBytes), err.Error())
	}
	return params, err
}

func HttpPostForm(urlStr string, body map[string]string, header map[string]string) (map[string]interface{}, error) {
	dataUrlVal := url.Values{}
	for key, val := range body {
		dataUrlVal.Add(key, val)
	}

	request, err := http.NewRequest("POST", urlStr, strings.NewReader(dataUrlVal.Encode()))
	if err != nil {
		return nil, fmt.Errorf("HTTP NewRequest POST %s error: %s", urlStr, err.Error())
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	// 添加额外的请求头
	for key, val := range header {
		request.Header.Set(key, val)
	}
	ctx, cancel := context.WithTimeout(request.Context(), time.Second*60)
	defer cancel()
	request = request.WithContext(ctx)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("HTTP POST %s error: %+v %s", urlStr, request, err.Error())
	}
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("HTTP POST RealAll [%s] error: %+v %s", urlStr, resp.Body, err.Error())
	}

	params := make(map[string]interface{})
	// json转map
	err = json.Unmarshal(respBytes, &params)
	if err != nil {
		return nil, fmt.Errorf("HTTP Unmarshal [%s] error: %s", string(respBytes), err.Error())
	}
	return params, err
}

func HttpGet(url string, header map[string]string) (cancelFunc context.CancelFunc, resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("HttpGet.NewRequest: %+v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	// 添加额外的请求头
	for key, val := range header {
		req.Header.Set(key, val)
	}
	ctx, cancel := context.WithTimeout(req.Context(), time.Second*60)
	req = req.WithContext(ctx)
	resp, err = http.DefaultClient.Do(req)
	return cancel, resp, err
}

func HttpResponseParseByJson(resp *http.Response) (map[string]interface{}, error) {
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	params := make(map[string]interface{})
	err = json.Unmarshal(b, &params)
	if err != nil {
		return nil, err
	}
	return params, err
}
