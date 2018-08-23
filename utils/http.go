package utils

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var cookie string

func HttpPost(url string, args map[string][]string) string {
	return httpQuery("POST", url, args)
}

func HttpGet(urls string) string {
	return httpQuery("GET", urls, url.Values{})
}

func httpQuery(method string, urls string, arguments map[string][]string) string {
	args := url.Values{}
	for k, v := range arguments {
		args.Set(k, v[0])
	}
	reqBody := ioutil.NopCloser(strings.NewReader(args.Encode()))
	request, _ := http.NewRequest(method, urls, reqBody)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	request.Header.Set("Cookie", cookie)
	request.Header.Set("Accept-Encoding", "identity")
	client := &http.Client{Timeout: time.Second * 2}
	resp, err := client.Do(request)
	if err != nil {
		return ""
	}

	c := resp.Header.Get("Cookie")
	if c != "" {
		cookie = c
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	result := string(respBody)
	return result
}

func HttpPostJson(urls string, json string) string {
	request, err := http.NewRequest("POST", urls, bytes.NewBuffer([]byte(json)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Cookie", cookie)
	request.Header.Set("Accept-Encoding", "identity")
	client := &http.Client{Timeout: time.Second * 2}
	resp, err := client.Do(request)
	if err != nil {
		return ""
	}

	c := resp.Header.Get("Cookie")
	if c != "" {
		cookie = c
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	result := string(respBody)
	return result
}
