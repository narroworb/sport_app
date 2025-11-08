package jwt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateAndParseToken(t *testing.T) {
	username := "some name"

	token, err := GenerateToken(username)
	if err != nil {
		t.Fatalf("error in generate token: %v", err)
	}

	usernameFromToken, err := ParseToken(token)
	if err != nil {
		t.Fatalf("error in parsing token: %v", err)
	}

	assert.Equal(t, username, usernameFromToken)
}

func TestParseBadToken(t *testing.T) {
	randomToken := "fsdjklfghjksdhg"

	_, err := ParseToken(randomToken)

	assert.Error(t, err)
}
