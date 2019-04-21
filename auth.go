package bankid

import (
	"encoding/json"
	"fmt"
)

type AuthResponse struct {
	AutoStartToken string `json:"autoStartToken"` // UUID, e.g "dbbee61c-357b-4fd8-b103-392eed10be7a"
	OrderRef       string `json:"orderRef"`       // UUID, e.g "131daac9-16c6-4618-beb0-365768f37288"
}

func Auth(env Environmenter, personalNumber string, userIP string) (*AuthResponse, error) {

	requestBody := Request{
		PersonalNumber: personalNumber,
		EndUserIP:      userIP,
	}

	req, err := env.NewRequest(AuthEndpoint, requestBody)
	if err != nil {
		return nil, err
	}

	client := env.NewClient() // A http.Client with a HTTP Mutal Authentication loaded

	rsp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer rsp.Body.Close()

	fmt.Printf("Auth response: %+v\n", rsp)

	authRsp := AuthResponse{}

	if err := json.NewDecoder(rsp.Body).Decode(&authRsp); err != nil {
		return nil, fmt.Errorf("could not decode authentication response: %s", err.Error())
	}

	fmt.Printf("%+v\n", authRsp)
	// TODO: Handle HTTP 4xx response

	return &authRsp, nil
}
