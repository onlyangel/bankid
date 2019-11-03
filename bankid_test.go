package bankid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvalidEnvironment(t *testing.T) {
	// Bad CA file
	_, err := NewEnvironment("BASE_URL", "INVALID_CA", "", "")
	assert.NotNil(t, err)

	// Bad Cert/Key files
	_, err = NewEnvironment("BASE_URL", "./CA/test.crt", "INVALID_RP_CERT", "INVALID_RP_KEY")
	assert.NotNil(t, err)

	// Wrong CA file
	_, err = NewEnvironment("BASE_URL", "./rp/bankid_rp_test.key", "./rp/bankid_rp_test.crt", "./rp/bankid_rp_test.key")
	assert.NotNil(t, err)
}

func TestValidEnvironment(t *testing.T) {
	// Cert Files OK
	_, err := NewEnvironment("BASE_URL", "./CA/test.crt", "./rp/bankid_rp_test.crt", "./rp/bankid_rp_test.key")
	assert.Nil(t, err)
}

func TestRequestsBad(t *testing.T) {
	env, err := NewEnvironment("BASE_URL:🤣", "./CA/test.crt", "./rp/bankid_rp_test.crt", "./rp/bankid_rp_test.key")
	assert.Nil(t, err)
	assert.NotNil(t, env)

	// Bad body
	var invalidBodyType chan int
	_, err = env.NewRequest("endpoint", invalidBodyType)
	assert.NotNil(t, err)

	// Bad schema
	_, err = env.NewRequest("endpoint", "")
	assert.NotNil(t, err)
}

func TestRequestsOK(t *testing.T) {
	env, err := NewEnvironment(ProductionBaseURL, "./CA/test.crt", "./rp/bankid_rp_test.crt", "./rp/bankid_rp_test.key")
	assert.Nil(t, err)
	assert.NotNil(t, env)

	// All A OK
	req, err := env.NewRequest("endpoint", "")
	assert.Nil(t, err)
	assert.NotNil(t, req)

	// Initate client
	client := env.NewClient()
	assert.NotNil(t, client)
}
