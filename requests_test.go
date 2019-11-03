package bankid

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testEnv struct {
	handler func(w http.ResponseWriter, r *http.Request) // Reset this one for every test
	server  *httptest.Server
}

func (t *testEnv) NewClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, t.server.Listener.Addr().String())
			},
			Dial: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 5 * time.Second,
			IdleConnTimeout:     90 * time.Second,
		},
		Timeout: 10 * time.Second,
	}
}
func (t *testEnv) NewRequest(endpoint string, body interface{}) (*http.Request, error) {
	t.server = httptest.NewServer(http.HandlerFunc(t.handler))

	fmt.Println("Test server at ", t.server.Listener.Addr().String())

	requestBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// fmt.Println("Request to: ", APIVersion+endpoint)
	bodyReader := strings.NewReader(string(requestBody))
	req, err := http.NewRequest("POST", "http://"+t.server.URL+"/"+APIVersion+endpoint, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "Application/json")
	return req, nil

}

type tt struct {
	name    string
	handler func(w http.ResponseWriter, r *http.Request)
	assert  func(resp interface{}, err error)
}

//
// Signing and Auth, veri similar test cases
//

func TestSignAuthCollect_v5(t *testing.T) {
	testFunctions := []*tt{
		&tt{
			name: "Expected: Successful, w/ correct response body",
			handler: func(w http.ResponseWriter, r *http.Request) {
				authRsp := Response{} // Empty but valid
				w.WriteHeader(200)
				json.NewEncoder(w).Encode(&authRsp)
			},
			assert: func(resp interface{}, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, resp)
			},
		},
		&tt{
			name: "Expected: Successful, w/ incorrect response body",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write(nil) // Empty and invalid
			},
			assert: func(resp interface{}, err error) {
				assert.NotNil(t, err)
				assert.Empty(t, resp)
			},
		},
		&tt{
			name: "Expected: Fail, w/ correct response body",
			handler: func(w http.ResponseWriter, r *http.Request) {
				errRsp := ErrorResponse{} // Empty but valid
				w.WriteHeader(400)
				json.NewEncoder(w).Encode(&errRsp)
			},
			assert: func(resp interface{}, err error) {
				assert.NotNil(t, err)
				assert.NotEmpty(t, err.Error())
				assert.Empty(t, resp)
			},
		},
		&tt{
			name: "Expected: Fail, w/ incorrect response body",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(400)
				w.Write(nil) // Empty and invalid
			},
			assert: func(resp interface{}, err error) {
				assert.NotNil(t, err)
				assert.NotEmpty(t, err.Error())
				assert.Empty(t, resp)
			},
		},
		&tt{ // Mostly f or test coverage
			name: "Expected: Fail, we don't handle 1xx/3xx messages",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(101)
				w.Write(nil) // Empty and invalid
			},
			assert: func(resp interface{}, err error) {
				assert.Nil(t, err)
				assert.Empty(t, resp)
			},
		},
	}

	env := &testEnv{}

	// Signing data
	for _, at := range testFunctions {
		env.handler = at.handler
		fmt.Printf("Sign: %s - ", at.name)
		resp, err := Sign(env, "198001010000", "127.0.0.1", "Hi User", "abc123")
		at.assert(resp, err)
	}

	// Authentication self
	for _, at := range testFunctions {
		env.handler = at.handler
		fmt.Printf("Auth: %s - ", at.name)
		resp, err := Auth(env, "198001010000", "127.0.0.1")
		at.assert(resp, err)
	}

	// Collecting status
	for _, at := range testFunctions {
		env.handler = at.handler
		fmt.Printf("Collect: %s - ", at.name)
		resp, err := Collect(env, "dbbee61c-357b-4fd8-b103-392eed10be7a")
		at.assert(resp, err)
	}

	env.server.Close()
}

func TestCancel_v5(t *testing.T) {
	testFunctions := []*tt{
		&tt{
			name: "Expected: Successful, w/ correct response body",
			handler: func(w http.ResponseWriter, r *http.Request) {
				authRsp := Response{} // Empty but valid
				w.WriteHeader(200)
				json.NewEncoder(w).Encode(&authRsp)
			},
			assert: func(_ interface{}, err error) {
				assert.Nil(t, err)
			},
		},
		&tt{
			name: "Expected: Successful, w/ incorrect response body",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write(nil) // Empty and invalid
			},
			assert: func(_ interface{}, err error) {
				assert.NotNil(t, err)
			},
		},
		&tt{
			name: "Expected: Fail, w/ correct response body",
			handler: func(w http.ResponseWriter, r *http.Request) {
				errRsp := ErrorResponse{} // Empty but valid
				w.WriteHeader(400)
				json.NewEncoder(w).Encode(&errRsp)
			},
			assert: func(_ interface{}, err error) {
				assert.NotNil(t, err)
			},
		},
	}

	env := &testEnv{}

	for _, at := range testFunctions {
		env.handler = at.handler
		fmt.Printf("Cancel: %s - ", at.name)
		err := Cancel(env, "dbbee61c-357b-4fd8-b103-392eed10be7a")
		at.assert(nil, err)
	}

	env.server.Close()
}

//
// Test invalid environment
//

type invalidEnv struct {
	testEnv
	request      *http.Request
	requestError error
}

func (t *invalidEnv) NewRequest(endpoint string, body interface{}) (*http.Request, error) {
	return t.request, t.requestError
}

func TestCallMethod(t *testing.T) {
	env := &invalidEnv{}

	env.request = &http.Request{}
	env.requestError = nil
	req, err := call("", env, nil, nil)
	assert.Nil(t, req)
	assert.NotNil(t, err)

	env.requestError = fmt.Errorf("fake invalid response")
	req, err = call("", env, nil, nil)
	assert.Nil(t, req)
	assert.NotNil(t, err)
}
