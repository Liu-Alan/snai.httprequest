package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func HttpRequestForm(method string, urlPath string, headerMap map[string]string, parms map[string]string) ([]byte, error) {
	if urlPath == "" {
		return []byte{}, errors.New("url is empty ")
	}

	var resParms io.Reader
	if parms != nil {
		parmsVal := url.Values{}
		for k, v := range parms {
			parmsVal.Add(k, v)
		}
		resParms = io.NopCloser(strings.NewReader(parmsVal.Encode()))
	} else {
		resParms = nil
	}

	request, err := http.NewRequest(method, urlPath, resParms)
	if err != nil {
		return []byte{}, errors.New("request error 101")
	}

	if headerMap != nil {
		for k, v := range headerMap {
			request.Header.Add(k, v)
		}
	}

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return []byte{}, errors.New("request error 102")
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return []byte{}, errors.New("response error 103")
	}
	defer client.CloseIdleConnections()

	return body, nil
}

func HttpRequestJson(method string, urlPath string, headerMap map[string]string, parms map[string]interface{}) ([]byte, error) {
	if urlPath == "" {
		return []byte{}, errors.New("url is empty ")
	}

	var resParms io.Reader
	if parms != nil {
		parmsByte, err := json.Marshal(parms)
		if err != nil {
			return []byte{}, errors.New("parms error")
		}
		resParms = io.NopCloser(bytes.NewReader(parmsByte))
	} else {
		resParms = nil
	}

	request, err := http.NewRequest(method, urlPath, resParms)
	if err != nil {
		return []byte{}, errors.New("request error 101")
	}

	if headerMap != nil {
		for k, v := range headerMap {
			request.Header.Add(k, v)
		}
	}

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return []byte{}, errors.New("request error 102")
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return []byte{}, errors.New("response error 103")
	}
	defer client.CloseIdleConnections()

	return body, nil
}

type ResValue struct {
	Code    int      `json:"code"`
	Msg     string   `json:"msg"`
	ResData *ResData `json:"data"`
}

type ResData struct {
	UserID  int    `json:"user_id"`
	Account string `json:"account"`
}

func main() {
	fmt.Println("form请求参数")
	start := time.Now()

	headerMap := make(map[string]string)
	headerMap["Accept"] = "application/json, text/plain, */*"
	headerMap["Accept-Encoding"] = "gzip, deflate, br, zstd"
	headerMap["Accept-Language"] = "zh-CN,zh;q=0.9"
	headerMap["Connection"] = "keep-alive"
	headerMap["Content-Type"] = "application/x-www-form-urlencoded"
	headerMap["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36"

	parms := make(map[string]string)
	parms["account"] = "snai"
	parms["password"] = "snai"

	response, err := HttpRequestForm("POST", "http://localhost:8080/user/login", headerMap, parms)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var resValue ResValue
	err = json.Unmarshal(response, &resValue)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("user_id:%d,account:%s,%v ms\n", resValue.ResData.UserID, resValue.ResData.Account, time.Since(start).Milliseconds())

	fmt.Println("json请求参数")
	start = time.Now()

	headerMap = make(map[string]string)
	headerMap["Accept"] = "application/json, text/plain, */*"
	headerMap["Accept-Encoding"] = "gzip, deflate, br, zstd"
	headerMap["Accept-Language"] = "zh-CN,zh;q=0.9"
	headerMap["Connection"] = "keep-alive"
	headerMap["Content-Type"] = "application/json"
	headerMap["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36"

	parmsJson := make(map[string]interface{})
	parmsJson["user_id"] = 10000
	parmsJson["password"] = "snai"

	response, err = HttpRequestJson("POST", "http://localhost:8080/user/login", headerMap, parmsJson)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = json.Unmarshal(response, &resValue)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("user_id:%d,account:%s,%v ms\n", resValue.ResData.UserID, resValue.ResData.Account, time.Since(start).Milliseconds())

	fmt.Scanln()
}
