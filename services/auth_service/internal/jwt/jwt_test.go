package jwt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateAndParseToken(t *testing.T) {
	username := "some name"
	var userID int64 = 4

	token, err := GenerateToken(username, userID)
	if err != nil {
		t.Fatalf("error in generate token: %v", err)
	}

	usernameFromToken, userIDFromTokentoken, err := ParseToken(token)
	if err != nil {
		t.Fatalf("error in parsing token: %v", err)
	}

	assert.Equal(t, username, usernameFromToken)
	assert.Equal(t, userID, userIDFromTokentoken)
}

func TestParseBadToken(t *testing.T) {
	randomToken := "fsdjklfghjksdhg"

	_, _, err := ParseToken(randomToken)

	assert.Error(t, err)
}
