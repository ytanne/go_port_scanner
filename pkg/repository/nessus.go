package repository

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/ytanne/go_nessus/pkg/entities"
)

type NessusClient struct {
	AccessKey string
	SecretKey string
	URL       string
	TR        *http.Transport
}

func NewNessusClient(AccessKey, SecretKey, URL string) *NessusClient {
	var tr *http.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &NessusClient{AccessKey: AccessKey, SecretKey: SecretKey, URL: URL, TR: tr}
}

func (n *NessusClient) ListScans() (*entities.ScanList, error) {
	client := &http.Client{Transport: n.TR}
	response, err := client.Get(n.URL + "scans")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var result entities.ScanList
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
