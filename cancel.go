package bankid

import (
	"encoding/json"
	"fmt"
)

func Cancel(env Environmenter, orderRef string) error {
	requestBody := Request{
		OrderRef: orderRef,
	}

	req, err := env.NewRequest(CancelEndpoint, requestBody)
	if err != nil {
		return err
	}

	client := env.NewClient() // A http.Client with a HTTP Mutal Authentication loaded

	rsp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer rsp.Body.Close()

	fmt.Printf("Cancel response: %+v\n", rsp)

	response := make(map[string]string)

	if err := json.NewDecoder(rsp.Body).Decode(&response); err != nil {
		return fmt.Errorf("could not cancel order: %s", err.Error())
	}

	fmt.Printf("Cancel response body: %+v\n", response)
	// TODO: Handle HTTP 4xx response

	return nil
}
