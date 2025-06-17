package util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)


func TestConnectivity(url string, timeout time.Duration) bool {
	client := http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:   true, 
			MaxIdleConnsPerHost: 512,  
		},
		Timeout: timeout, 
	}
	req, err := http.NewRequest("GET", url, bytes.NewReader([]byte("")))
	if err != nil {
		return false
	}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	return true
}

func HttpGet(url string, timeout time.Duration) (string, error) {
	client := http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:   true, 
			MaxIdleConnsPerHost: 512,  
		},
		Timeout: timeout, 
	}
	req, err := http.NewRequest("GET", url, bytes.NewReader([]byte("")))
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	
	defer resp.Body.Close()
	respContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(respContent), err
}

func HttpPostJson(url string, requestBody []byte) (string, error) {
	response, er := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if er != nil {
		return "", er
	}
	defer response.Body.Close()
	body, er2 := ioutil.ReadAll(response.Body)
	if er2 != nil {
		return "", er2
	}
	return string(body), nil
}

func HttpPost(url string, requestBody []byte, contentType string) (string, error) {
	client := http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:   true, 
			MaxIdleConnsPerHost: 512,  
		},
		Timeout: time.Second * 5, 
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(requestBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/octet-stream")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	
	fmt.Println(resp.Header)
	defer resp.Body.Close()
	respContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(respContent), err
}
