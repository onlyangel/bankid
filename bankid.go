package bankid

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// API Constants
const (
	ProductionBaseURL string = "https://appapi2.bankid.com"
	TestBaseURL       string = "https://appapi2.test.bankid.com"
	APIVersion        string = "/rp/v5"
	AuthEndpoint      string = "/auth"
	SignEndpoint      string = "/sign"
	CollectEndpoint   string = "/collect"
	CancelEndpoint    string = "/cancel"
)

// HintCodes - for Pending and Failed statuses
const (
	PendOutstandingTransaction = "outstandingTransaction"
	PendNoClient               = "noClient"
	PendStarted                = "started"
	PendUserSign               = "userSign"
	// PendUnknown

	FailExpiredTransaction = "expiredTransaction"
	FailCertificateErr     = "certificateErr"
	FailUserCancel         = "userCancel"
	FailCancelled          = "cancelled"
	FailStartFailed        = "startFailed"
	// FailUnknown
)

// Environmenter  ¯\_(ツ)_/¯
// Helps setup requests to the BankID API
type Environmenter interface {
	NewClient() *http.Client
	NewRequest(endpoint string, body interface{}) (*http.Request, error)
}

type environment struct {
	baseURL      string
	clientConfig *tls.Config
}

// Request - A basic BankID request contain one or more of the variables below
type Request struct {
	OrderRef           string `json:"orderRef,omitempty"`
	EndUserIP          string `json:"endUserIp,omitempty"`
	PersonalNumber     string `json:"personalNumber,omitempty"`
	UserVisibleData    string `json:"userVisibleData,omitempty"`
	UserNonVisibleData string `json:"userNonVisibleData,omitempty"`
}

func NewEnvironment(baseURL string, caPath string, rpCertPath string, rpKeyPath string) (*environment, error) {
	ca, err := ioutil.ReadFile(caPath)
	if err != nil {
		return nil, fmt.Errorf("Could not load CA Certificate: %s", err.Error())
	}

	rpCert, err := tls.LoadX509KeyPair(rpCertPath, rpKeyPath)
	if err != nil {
		return nil, fmt.Errorf("Could not load RP Keypair: %s", err.Error())
	}

	caPool := x509.NewCertPool()

	if caPool.AppendCertsFromPEM(ca) == false {
		return nil, fmt.Errorf("Could not append CA Certificate to pool. Invalid certificate?")
	}

	clientCfg := tls.Config{
		Certificates: []tls.Certificate{rpCert},
		ClientCAs:    caPool,
		RootCAs:      caPool,
		// InsecureSkipVerify: true, // For some reason is BankID not using a proper domain certificate
	}
	return &environment{
		baseURL:      baseURL,
		clientConfig: &clientCfg,
	}, nil
}

func (e *environment) NewRequest(endpoint string, body interface{}) (*http.Request, error) {
	requestBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	fmt.Printf("New request to >> %s << with body: %s\n", endpoint, string(requestBody))

	bodyReader := strings.NewReader(string(requestBody))
	req, err := http.NewRequest("POST", e.baseURL+APIVersion+endpoint, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "Application/json")
	return req, nil
}

func (e *environment) NewClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: e.clientConfig,
		},
	}
}
