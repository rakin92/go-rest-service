// Package web is an wrapper around the net/http package to do basic http requests
package web

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/rakin92/go-rest-service/pkg/logger"
)

var c *http.Client

func init() {
	c = NewClient(10)
}

func NewClient(timeout int) *http.Client {
	t := &http.Transport{
		ForceAttemptHTTP2:     true,
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := &http.Client{
		Transport: t,
		Timeout:   time.Duration(timeout) * time.Second,
	}
	return client
}

func handleRequest(req *http.Request) ([]byte, error) {
	response, err := c.Do(req)
	if err != nil {
		logger.Error(&err, "Error sending request")
		return nil, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Error(&err, "Couldn't parse response body")
		return nil, err
	}

	return body, nil
}

// Get makes a http get request to the given url
func Get(fullURL string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		logger.Error(&err, "Error trying to make GET request")
		return nil, err
	}

	return handleRequest(req)
}

// Post makes a http post request to the given url and data
func Post(fullURL string, data any) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		logger.Error(&err, "Error marshalling request POST data")
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fullURL, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Error(&err, "Error trying to make POST request")
		return nil, err
	}

	return handleRequest(req)
}

// Put makes a http put request to the given url and data
func Put(fullURL string, data any) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		logger.Error(&err, "Error marshalling request PUT data")
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, fullURL, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Error(&err, "Error trying to make PUT request")
		return nil, err
	}

	return handleRequest(req)
}

// Patch makes a http patch request to the given url and data
func Patch(fullURL string, data any) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		logger.Error(&err, "Error marshalling request PATCH data")
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPatch, fullURL, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Error(&err, "Error trying to make PATCH request")
		return nil, err
	}

	return handleRequest(req)
}
