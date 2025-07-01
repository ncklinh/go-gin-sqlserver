package token

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker("12345678901234567890123456789012")
	require.NoError(t, err)
	require.NotNil(t, maker)

	maker, err = NewJWTMaker("short-key")
	require.Error(t, err)
	require.Nil(t, maker)
}

func TestCreateToken(t *testing.T) {
	maker, err := NewJWTMaker("12345678901234567890123456789012")
	require.NoError(t, err)

	username := "testuser"
	role := "admin"
	duration := time.Minute

	token, err := maker.CreateToken(username, role, duration, TokenTypeAccessToken)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.Len(t, strings.Split(token, "."), 3)
}

func TestVerifyToken_InvalidOrExpired(t *testing.T) {
	maker, err := NewJWTMaker("12345678901234567890123456789012")
	require.NoError(t, err)

	// Invalid token
	payload, err := maker.VerifyToken("invalid.token.here", TokenTypeAccessToken)
	require.Error(t, err)
	require.Nil(t, payload)
	require.Equal(t, ErrInvalidToken, err)

	// Expired token
	token, err := maker.CreateToken("testuser", "admin", -time.Minute, TokenTypeAccessToken)
	require.NoError(t, err)
	payload, err = maker.VerifyToken(token, TokenTypeAccessToken)
	require.Error(t, err)
	require.Nil(t, payload)
	require.Equal(t, ErrExpiredToken, err)
}
