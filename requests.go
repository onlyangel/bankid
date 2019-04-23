package bankid

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// Sign - besides personalNumber and userIP inputs as
// the Auth() call we also need data to sign.
//
// From documentation:
// User Visible data
//  The text to be displayed and signed. String.
//  The text can be formatted using CR, LF and CRLF for new lines.
//  The text must be encoded as UTF-8 and then base64 encoded.
//  1 to 40'000 characters after base64 encoding.
//
// User Non-visible data
//  Data not displayed to the user. String.
//  The value must be base64-encoded.
//  1 to 200'000 characters after base64-encoding.
//
// -------------------
// Because of base64-encoding adds 33% overhead, don't send more than
//  30'000 bytes of Visible data and
// 150'000 bytes of Non-visible data
//
// The Sign() method will base64-encode both the UserVisible and UserNonVisible data.
// Choose whichever line ending character you need.
func Sign(env Environmenter, personalNumber string, userIP string, userVisible string, userNonVisible string) (*Response, error) {

	// Base64 encode with padding
	if userVisible != "" {
		userVisible = base64.StdEncoding.EncodeToString([]byte(userVisible))
	}

	if userNonVisible != "" {
		userNonVisible = base64.StdEncoding.EncodeToString([]byte(userNonVisible))
	}

	requestBody := Request{
		PersonalNumber:     personalNumber,
		EndUserIP:          userIP,
		UserVisibleData:    userVisible,
		UserNonVisibleData: userNonVisible,
	}
	return call(SignEndpoint, env, &requestBody)
}

// Auth - verify a users' identitys
func Auth(env Environmenter, personalNumber string, userIP string) (*Response, error) {
	requestBody := Request{
		PersonalNumber: personalNumber,
		EndUserIP:      userIP,
	}
	return call(AuthEndpoint, env, &requestBody)
}

// Cancel -
func Cancel(env Environmenter, orderRef string) (*Response, error) {
	requestBody := Request{
		OrderRef: orderRef,
	}
	return call(CancelEndpoint, env, &requestBody)
}

func call(endpoint string, env Environmenter, requestBody *Request) (*Response, error) {

	req, err := env.NewRequest(endpoint, requestBody)
	if err != nil {
		return nil, err
	}

	client := env.NewClient() // A http.Client with a HTTP Mutal Authentication loaded

	rsp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer rsp.Body.Close()

	authRsp := Response{}

	if err := json.NewDecoder(rsp.Body).Decode(&authRsp); err != nil {
		return nil, fmt.Errorf("could not decode response: %s", err.Error())
	}
	// TODO: Handle HTTP 4xx response

	return &authRsp, nil
}
