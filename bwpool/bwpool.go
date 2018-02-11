package bwpool

import (
	"io/ioutil"
	"net/http"
	"time"
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
	"net/url"
	"strconv"
	"strings"

	"github.com/json-iterator/go"
)

// Why isn't this part of the stdlib?!
// Oh right: https://github.com/json-iterator/go/issues/231
var json = jsoniter.ConfigDefault

type PoolConfig struct {
	Username   string `json:"Username"`
	PublicKey  string `json:"PublicKey"`
	PrivateKey string `json:"PrivateKey"`
}

func fetchData(url string, postData string, timeout time.Duration) (*[]byte, error) {
	client := http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(postData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "bwpool_exporter")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &body, err
}

func workerValues(username string, publicKey string, nonce string, signature string) url.Values {
	v := url.Values{}
	v.Set("userName", username)
	v.Set("key", publicKey)
	v.Set("nonce", nonce)
	v.Set("signature", signature)
	v.Set("pageSize", "100")
	return v
}

func newWorkerResponse(body *[]byte) (*WorkerResponse, error) {
	wr := WorkerResponse{}

	err := json.Unmarshal(*body, &wr)
	if err != nil {
		return nil, err
	}

	return &wr, nil
}

type BWClient struct {
	urlWorkers string
	username   string
	publicKey  string
	privateKey string
	timeout    time.Duration
}

func (c *BWClient) FetchWorkers() (*PoolData, error) {
	nonce := strconv.FormatInt(time.Now().Unix(), 10)
	signature := getSignature(c.username, c.publicKey, c.privateKey, nonce)
	values := workerValues(c.username, c.publicKey, nonce, signature)

	body, err := fetchData(c.urlWorkers, values.Encode(), c.timeout)
	if err != nil {
		return nil, err
	}

	response, err := newWorkerResponse(body)
	if err != nil {
		return nil, err
	}

	workers := Workers{}
	for _, worker := range response.Workers {
		workers[worker.Name] = worker
	}
	pd := &PoolData{
		Workers: workers,
	}

	return pd, nil
}

func NewClient(poolConfig *PoolConfig, timeout time.Duration) *BWClient {
	return &BWClient{
		urlWorkers: "https://ltc.bw.com/api/workers",
		username:   poolConfig.Username,
		publicKey:  poolConfig.PublicKey,
		privateKey: poolConfig.PrivateKey,
		timeout:    timeout,
	}
}

func getSignature(username string, publicKey string, privateKey string, nonce string) string {
	message := username + publicKey + nonce
	signed := hmacSign([]byte(message), []byte(privateKey))
	encoded := hex.EncodeToString(signed)
	return encoded
}

func hmacSign(message []byte, key []byte) []byte {
	mac := hmac.New(md5.New, key)
	mac.Write(message)
	response := mac.Sum(nil)
	return response
}

func ReadConfig(path string) (*PoolConfig, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	pc := &PoolConfig{}
	err = json.Unmarshal(data, pc)
	if err != nil {
		return nil, err
	}

	return pc, nil
}
