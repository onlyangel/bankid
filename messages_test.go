package bankid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLanguageCases(t *testing.T) {

	validCases := []string{
		"en", "EN", "eN",
		"se", "SE", "sE",
	}

	for _, c := range validCases {
		_, err := NewMessages(c)
		assert.Nil(t, err)
	}
}

func TestUnknownLanguages(t *testing.T) {
	_, err := NewMessages("es") // No spanish support yet
	assert.NotNil(t, err)
}

func TestCorrectLanguages(t *testing.T) {
	// English
	en, err := NewMessages("en")
	assert.Nil(t, err)
	assert.Equal(t, messages_EN[RFA1], en.Msg(RFA1))

	// Swedish
	se, err := NewMessages("se")
	assert.Nil(t, err)
	assert.Equal(t, messages_SE[RFA1], se.Msg(RFA1))
}

func TestInvalidMessagesReference(t *testing.T) {
	se, err := NewMessages("se") // No spanish support yet
	assert.Nil(t, err)
	assert.Equal(t, "", se.Msg("INVALID_REFERENCE"))
}
