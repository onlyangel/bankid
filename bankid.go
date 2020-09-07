package bankid

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	pkcs12 "software.sslmate.com/src/go-pkcs12"
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

func NewEnvironmentP12(baseURL string, caPath string, rpP12Path string) (Environmenter, error) {
	ca, err := ioutil.ReadFile(caPath)
	if err != nil {
		return nil, fmt.Errorf("could not load CA Certificate: %s", err.Error())
	}

	p12bytres, err :=  ioutil.ReadFile(rpP12Path)
	if err != nil {
		return nil, fmt.Errorf("could not load p12 certificate: %s", err.Error())
	}

	blocks, err := pkcs12.ToPEM(p12bytres, "qwerty123")
	if err != nil {
		return nil, fmt.Errorf("could not converto to pem: %s", err.Error())
	}

	var pemData []byte
	for _, b := range blocks {
		pemData = append(pemData, pem.EncodeToMemory(b)...)
	}

	// then use PEM data for tls to construct tls certificate:
	rpCert, err := tls.X509KeyPair(pemData, pemData)
	if err != nil {
		return nil, fmt.Errorf("could not load RP Keypair: %s", err.Error())
	}

	caPool := x509.NewCertPool()

	if caPool.AppendCertsFromPEM(ca) == false {
		return nil, fmt.Errorf("could not append CA Certificate to pool. Invalid certificate?")
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


// NewEnvironment - sets up the certificates and URLs needed to identify ourselves with the BankID service
func NewEnvironment(baseURL string, caPath string, rpCertPath string, rpKeyPath string) (Environmenter, error) {
	ca, err := ioutil.ReadFile(caPath)
	if err != nil {
		return nil, fmt.Errorf("could not load CA Certificate: %s", err.Error())
	}

	rpCert, err := tls.LoadX509KeyPair(rpCertPath, rpKeyPath)
	if err != nil {
		return nil, fmt.Errorf("could not load RP Keypair: %s", err.Error())
	}

	caPool := x509.NewCertPool()

	if caPool.AppendCertsFromPEM(ca) == false {
		return nil, fmt.Errorf("could not append CA Certificate to pool. Invalid certificate?")
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

// NewRequest - helper function to bake a request
func (e *environment) NewRequest(endpoint string, body interface{}) (*http.Request, error) {
	requestBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	bodyReader := strings.NewReader(string(requestBody))
	req, err := http.NewRequest("POST", e.baseURL+APIVersion+endpoint, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "Application/json")
	return req, nil
}

// NewRequest - helper function to bake a new http.Client with our TLS Confnig
func (e *environment) NewClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSClientConfig:     e.clientConfig,
			TLSHandshakeTimeout: 5 * time.Second,
			IdleConnTimeout:     90 * time.Second,
		},
		Timeout: 10 * time.Second,
	}
}
