package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/facebookgo/httpcontrol"
)

var httpClient *http.Client

func init() {
	httpClient = &http.Client{
		Transport: &httpcontrol.Transport{
			RequestTimeout: 10 * time.Second,
			MaxTries:       3,
		},
	}
}

func gethttpClient(url string) *http.Client {
	return &http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{}},
		Timeout:   time.Duration(10 * time.Second),
	}
}

func request(path, method string, body io.Reader) ([]byte, error) {

	req, errRequest := http.NewRequest(method, path, body)
	if errRequest != nil {
		return nil, errRequest
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "close")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp == nil {
		return nil, fmt.Errorf("Invalid response from Bitbucket. Please Check Bitbucket, err:%s", err)
	}
	if resp.StatusCode >= 300 {
		log.Errorf("Response Status:%s", resp.Status)
		log.Errorf("Response Headers:%s", resp.Header)
	}

	bodyResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error with ioutil.ReadAll %s", err.Error())
	}
	return bodyResp, nil
}
