package bankid

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Use this to parse the BankID API response, see stdResponseParser
type responseParser func(*http.Response) (interface{}, error)

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

	output := &Response{}
	rsp, err := call(SignEndpoint, env, &requestBody, stdResponseParser)
	if err == nil {
		output = rsp.(*Response)
	}
	return output, err
}

// Auth - verify a users identity
func Auth(env Environmenter, personalNumber string, userIP string) (*Response, error) {
	requestBody := Request{
		PersonalNumber: personalNumber,
		EndUserIP:      userIP,
	}
	output := &Response{}
	rsp, err := call(AuthEndpoint, env, &requestBody, stdResponseParser)
	if err == nil {
		output = rsp.(*Response)
	}
	return output, err
}

func Collect(env Environmenter, orderRef string) (*CollectResponse, error) {
	requestBody := Request{
		OrderRef: orderRef,
	}

	output := &CollectResponse{}
	rsp, err := call(CollectEndpoint, env, &requestBody, collectParser)
	if err == nil {
		output = rsp.(*CollectResponse)
	}
	return output, err
}

// Cancel -
func Cancel(env Environmenter, orderRef string) (*Response, error) {
	requestBody := Request{
		OrderRef: orderRef,
	}
	output := &Response{}
	rsp, err := call(CancelEndpoint, env, &requestBody, stdResponseParser)
	if err == nil {
		output = rsp.(*Response)
	}
	return output, err
}

func call(endpoint string, env Environmenter, requestBody *Request, rspParser responseParser) (interface{}, error) {

	req, err := env.NewRequest(endpoint, requestBody)
	if err != nil {
		return nil, err
	}

	client := env.NewClient() // A http.Client with a HTTP Mutal Authentication loaded

	rsp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return rspParser(rsp)
}

func stdResponseParser(rsp *http.Response) (interface{}, error) {
	defer rsp.Body.Close()

	log.Printf("Response code: %s\n", rsp.Status)

	// OK
	if rsp.StatusCode >= 200 && rsp.StatusCode < 400 {
		authRsp := Response{}
		err := json.NewDecoder(rsp.Body).Decode(&authRsp)
		if err != nil {
			panic(fmt.Sprintf("failed to parse successful response: %s", err.Error()))
		}
		return &authRsp, nil
	}

	// Fail
	if rsp.StatusCode >= 400 {
		errRsp := ErrorResponse{}
		err := json.NewDecoder(rsp.Body).Decode(&errRsp)
		if err != nil {
			panic(fmt.Sprintf("failed to parse fail response: %s", err.Error()))
		}
		return nil, errRsp // A bit unorthodox but ErrorResponse is a proper error type
	}

	// We don't care about HTTP 1xx messages
	return nil, nil
}

func collectParser(rsp *http.Response) (interface{}, error) {
	defer rsp.Body.Close()

	log.Printf("Response code: %s\n", rsp.Status)

	// OK
	if rsp.StatusCode >= 200 && rsp.StatusCode < 400 {
		collectRsp := CollectResponse{}
		err := json.NewDecoder(rsp.Body).Decode(&collectRsp)
		if err != nil {
			panic(fmt.Sprintf("failed to parse successful response: %s", err.Error()))
		}
		return &collectRsp, nil
	}

	// Fail
	if rsp.StatusCode >= 400 {
		errRsp := ErrorResponse{}
		err := json.NewDecoder(rsp.Body).Decode(&errRsp)
		if err != nil {
			panic(fmt.Sprintf("failed to parse fail response: %s", err.Error()))
		}
		return nil, errRsp // A bit unorthodox but ErrorResponse is a proper error type
	}
	// We don't care about HTTP 1xx messages
	return nil, nil
}
