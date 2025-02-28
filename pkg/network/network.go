package network

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// 최적화된 Transport 설정
var customTransport = &http.Transport{
	MaxIdleConns:        100,              // 최대 유휴 연결 수
	MaxConnsPerHost:     10,               // 호스트당 최대 연결 수
	IdleConnTimeout:     90 * time.Second, // 유휴 연결 타임아웃
	TLSHandshakeTimeout: 10 * time.Second, // TLS 핸드셰이크 타임아웃
}

// 재사용 가능한 HTTP Client
var httpClient = &http.Client{
	Transport: customTransport,
	Timeout:   10 * time.Second, // 요청 타임아웃
}

// 기본 헤더와 쿠키를 처리하는 함수들
var defaultHeaders = map[string]string{
	"Accept":          "application/json, text/plain, */*",
	"Accept-Encoding": "gzip",
	"Accept-Language": "ko-KR,ko;q=0.9,en-US;q=0.8,en;q=0.7",
	"Connection":      "keep-alive",
	"User-Agent":      "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
}

func addDefaultHeader(req *http.Request) {
	for key, value := range defaultHeaders {
		req.Header.Set(key, value)
	}
}

func addCookies(req *http.Request, cookies map[string]string) {
	for name, value := range cookies {
		req.AddCookie(&http.Cookie{Name: name, Value: value})
	}
}

// Gzip 응답 처리 함수
func decodeGzipBody(body io.ReadCloser) ([]byte, error) {
	defer body.Close()

	gzipReader, err := gzip.NewReader(body)
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()

	// Gzip 압축 해제된 데이터를 읽음
	return io.ReadAll(gzipReader)
}

// GET 요청 함수
func GetRequest(url string, cookies map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// 기본 헤더 및 쿠키 추가
	addDefaultHeader(req)
	if cookies != nil {
		addCookies(req, cookies)
	}

	// 요청 실행
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Gzip 응답인지 확인 후 처리
	if res.Header.Get("Content-Encoding") == "gzip" {
		return decodeGzipBody(res.Body)
	}

	// 응답 본문 읽기
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// POST 요청 함수 (application/x-www-form-urlencoded)
func PostRequest(url string, formData url.Values, cookies map[string]string) ([]byte, error) {
	payload := strings.NewReader(formData.Encode())

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}

	// 기본 헤더 및 쿠키 추가
	addDefaultHeader(req)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	addCookies(req, cookies)

	// 요청 실행
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Gzip 응답인지 확인 후 처리
	if res.Header.Get("Content-Encoding") == "gzip" {
		return decodeGzipBody(res.Body)
	}

	// 응답 본문 읽기
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// POST 요청 함수 (application/json)
func PostJSONRequest(url string, payload string, cookies map[string]string) ([]byte, *http.Response, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		return nil, nil, err
	}

	// 기본 헤더 및 쿠키 추가
	addDefaultHeader(req)
	req.Header.Set("Content-Type", "application/json")
	addCookies(req, cookies)

	// 요청 실행
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	// Gzip 응답인지 확인 후 처리
	if res.Header.Get("Content-Encoding") == "gzip" {
		decodeBody, err := decodeGzipBody(res.Body)
		if err != nil {
			return nil, nil, err
		}
		return decodeBody, res, nil
	}

	// 응답 본문 읽기
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}

	return body, res, nil
}
