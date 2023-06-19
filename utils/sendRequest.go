package utils

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html/charset"
	"golang.org/x/net/proxy"
)

type SendRequest struct {
	client    *http.Client
	retryNum  int // 重试次数
	retryTime int // 重试间隔时间
	Headers   http.Header
	boundary  string
}

func NewSendRequest(headers http.Header, boundary string) *SendRequest {
	if headers == nil {
		headers = http.Header{}
		headers.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	return &SendRequest{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		retryTime: 10,
		Headers:   headers,
		boundary:  boundary,
	}
}

func (s *SendRequest) SetHeaders(headers map[string]string) {
	// 设置请求头
	for key, value := range headers {
		s.Headers.Set(key, value)
	}
}

func (s *SendRequest) SetProxy(proxyAddr string, t string) {
	// 创建代理 URL
	proxyURL, err := url.Parse(proxyAddr)
	if err != nil {
		fmt.Println("Failed to parse proxy URL:", err)
		return
	}
	var transport *http.Transport
	if strings.Contains(t, "socks") {
		dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
		if err != nil {
			fmt.Println("Failed to create proxy dialer:", err)
			return
		}

		// 创建自定义的 HTTP 客户端，使用 SOCKS5 代理进行请求
		transport = &http.Transport{
			Dial: dialer.Dial,
		}

		// proxy := func(_ *http.Request) (*url.URL, error) {
		// 	return url.Parse(proxyAddr)
		// }
		// transport = &http.Transport{
		// 	Proxy: proxy,
		// }
	} else {
		// 创建自定义的 Transport
		transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	s.client.Transport = transport
}

func (s *SendRequest) send(method string, url string, param url.Values) ([]byte, *http.Response, error) {
	reqBody := strings.NewReader(param.Encode())
	// 设置请求参数
	if s.boundary != "" {
		// 创建请求体
		buf := &bytes.Buffer{}
		writer := multipart.NewWriter(buf)
		// 设置分割符号（boundary）
		writer.SetBoundary(s.boundary)
		// 添加表单字段到请求体
		for key, value := range param {
			for _, v := range value {
				_ = writer.WriteField(key, v)
			}
		}
		// 关闭 multipart.Writer，以写入结尾标识符
		_ = writer.Close()
		reqBody = strings.NewReader(buf.String())
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, nil, err
	}

	// 设置请求头
	req.Header = s.Headers

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode != http.StatusOK {
		// return nil, nil, errors.New(fmt.Sprintf("状态码：%d，内容：%s", resp.StatusCode, string(body)))
		return body, resp, errors.New(fmt.Sprintf("状态码：%d", resp.StatusCode))
	}

	// 从响应头中获取字符集编码
	contentType := resp.Header.Get("Content-Type")
	charsetReader, err := charset.NewReader(strings.NewReader(string(body)), contentType)

	bodyString, err := ioutil.ReadAll(charsetReader)
	if err != nil {
		fmt.Println("字符集解码失败:", err)
		return nil, nil, err
	}

	return bodyString, resp, nil
}

func (s *SendRequest) Post(url string, param url.Values) ([]byte, *http.Response, error) {
	return s.send("POST", url, param)
}

func (s *SendRequest) RepeatSend(method string, url string, param url.Values) ([]byte, *http.Response, error) {
	// fmt.Printf("[%s][%d]请求地址：%s\n", uuid, num, url)
	// fmt.Printf("[%s][%d]本次发送：%v\n", uuid, num, param)
	var (
		result []byte
		resp   *http.Response
		err    error
	)
	switch method {
	case "GET":
		result, resp, err = s.Get(url)
	case "POST":
		result, resp, err = s.Post(url, param)
	}

	if err != nil && s.retryNum < 5 {
		time.Sleep(time.Second * time.Duration(s.retryTime))
		// fmt.Printf("[%s][%d]请求返回错误：%v\n", uuid, num, err)
		s.retryNum++
		s.retryNum = s.retryNum * 2
		// 进行重发
		return s.RepeatSend(method, url, param)
	}
	return result, resp, err
}

func (s *SendRequest) Get(url string) ([]byte, *http.Response, error) {
	return s.send("GET", url, nil)
}
