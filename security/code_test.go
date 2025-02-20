package security_test

import (
	"testing"

	"github.com/ngtrvu/zen-go/security"
	"github.com/stretchr/testify/assert"
)

func TestGenerateOTPNumber(t *testing.T) {
	s := security.GenerateOTPNumber(8)
	assert.NotEqual(t, "", s)
	assert.Equal(t, 8, len(s))
}

func TestGenerateSecureToken(t *testing.T) {
	s := security.GenerateSecureToken(8)
	assert.NotEqual(t, "", s)
	assert.Equal(t, 16, len(s))
}
