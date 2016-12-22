package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/facebookgo/httpcontrol"
)

var sslInsecureSkipVerify bool

func isHTTPS(url string) bool {
	return strings.HasPrefix(url, "https")
}

func getHTTPClient(url string) *http.Client {
	var tr *http.Transport
	if isHTTPS(url) {
		tlsConfig := getTLSConfig()
		tr = &http.Transport{TLSClientConfig: tlsConfig}
	} else {
		tr = &http.Transport{}
	}

	timeout := time.Duration(10 * time.Second)
	return &http.Client{
		Transport: tr,
		Timeout:   timeout,
	}
}

func getTLSConfig() *tls.Config {
	return &tls.Config{}
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// HTTPClient is HTTClient or testHTTPClient for tests
var HTTPClient httpClient

func reqWant(path, method string, jsonStr []byte) ([]byte, error) {

	var req *http.Request
	if jsonStr != nil {
		req, _ = http.NewRequest(method, path, bytes.NewReader(jsonStr))
	} else {
		req, _ = http.NewRequest(method, path, nil)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "close")

	if HTTPClient == nil {
		HTTPClient = &http.Client{
			Transport: &httpcontrol.Transport{
				RequestTimeout: 10 * time.Second,
				MaxTries:       3,
			},
		}
	}
	resp, err := HTTPClient.Do(req)

	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if resp == nil {
		return []byte{}, fmt.Errorf("Invalid response from Bitbucket. Please Check Bitbucket, err:%s", err)
	}
	if resp.StatusCode >= 300 {
		log.Errorf("Request Body:%s", string(jsonStr))
		log.Errorf("Response Status:%s", resp.Status)
		log.Errorf("Response Headers:%s", resp.Header)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error with ioutil.ReadAll %s", err.Error())
	}
	return body, nil
}
