package utils

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

func ParseCookieString(rawCookies string) (*http.Request, error) {
	rawRequest := fmt.Sprintf("GET / HTTP/1.0\r\nCookie: %s\r\n\r\n", rawCookies)
	req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(rawRequest)))
	return req, err
}

// var utilsBytesSeed []byte = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890")

func init() {
	// 保证每次生成的随机数不一样
	rand.Seed(time.Now().UnixNano())
}

// 方法二
func RandStr2(n int) string {
	result := make([]byte, n/2)
	rand.Read(result)
	return hex.EncodeToString(result)
}

func RandomInt(min, max int) int {
	// 设置随机种子
	rand.Seed(time.Now().UnixNano())

	// 生成随机数
	return rand.Intn(max-min+1) + min
}

// map转url.Values
func MapToUrlValue(params map[string]string) url.Values {
	values := url.Values{}
	for key, value := range params {
		values.Add(key, value)
	}
	return values
}

func IsStructEmpty(s interface{}) bool {
	// 获取结构体的反射类型
	t := reflect.TypeOf(s)

	// 如果不是结构体类型，则返回 false
	if t.Kind() != reflect.Struct {
		return false
	}

	// 遍历结构体的字段
	for i := 0; i < t.NumField(); i++ {
		// 获取字段的值
		fieldValue := reflect.ValueOf(s).Field(i)

		// 如果字段的值不是零值，则结构体不为空
		if !fieldValue.IsZero() {
			return false
		}
	}

	// 结构体为空
	return true
}
