package bankid

import (
	"encoding/json"
	"fmt"
)

// Collect response statuses
const (
	OrderPending  = "pending"
	OrderFailed   = "failed"
	OrderComplete = "complete"
)

type CollectResponse struct {
	OrderRef       string      `json:"orderRef"`
	Status         string      `json:"status"`
	HintCode       string      `json:"hintCode,omitempty"`       // Pending and Failed orders only
	CompletionData *Completion `json:"completionData,omitempty"` // Complete orders only
}

type Completion struct {
	User         User   `json:"user"`
	Device       Device `json:"device"`
	Cert         Cert   `json:"cert"`
	Signature    string `json:"signature"` // base64 encoded signature, see https://www.bankid.com/bankid-i-dina-tjanster/rp-info
	OCSPResponse string `json:"ocspResponse"`
}

type User struct {
	PersonalNumber string `json:"personalNumber"` // e.g "197001010000"
	Name           string `json:"name"`
	GivenName      string `json:"giveName"`
	Surname        string `json:"surname"`
}

type Device struct {
	IPAddress string `json:"ipAddress"` // e.g "192.168.0.1"
}

type Cert struct {
	NotBefore string `json:"notBefore"` // e.g "1502983274000" UNIX Epoch in ms
	NotAfter  string `json:"notAfter"`
}

func Collect(env Environmenter, orderRef string) (*CollectResponse, error) {
	requestBody := Request{
		OrderRef: orderRef,
	}

	req, err := env.NewRequest(CollectEndpoint, requestBody)
	if err != nil {
		return nil, err
	}

	client := env.NewClient() // A http.Client with a HTTP Mutal Authentication loaded

	rsp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer rsp.Body.Close()

	fmt.Printf("Collect response: %+v\n", rsp)

	response := CollectResponse{}

	if err := json.NewDecoder(rsp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("could not decode Collect response: %s", err.Error())
	}

	fmt.Printf("Collect response body: %+v\n", response)
	// TODO: Handle HTTP 4xx response

	return &response, nil

}
