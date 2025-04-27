package tokens

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Tokens struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int64  `json:"expires_in"`
	RefreshToken     string `json:"refresh_token"`
	RefreshExpiresIn int64  `json:"refresh_expires_in"`
	TokenType        string `json:"token_type"`
}

type claims struct {
	UserID int64 `json:"uid"`
	jwt.RegisteredClaims
}

type Manager struct {
	secret            []byte
	accessTTLSeconds  int64
	refreshTTLSeconds int64
}

func NewManager(secret []byte, accessTTL, refreshTTL int64) *Manager {
	return &Manager{secret, accessTTL, refreshTTL}
}

func (m *Manager) Generate(userID int64) (*Tokens, error) {
	now := time.Now()

	accessToken, err := m.signedToken(userID, now.Add(time.Duration(m.accessTTLSeconds)*time.Second))
	if err != nil {
		return nil, err
	}
	refreshToken, err := m.signedToken(userID, now.Add(time.Duration(m.refreshTTLSeconds)*time.Second))
	if err != nil {
		return nil, err
	}

	return &Tokens{
		AccessToken:      accessToken,
		ExpiresIn:        m.accessTTLSeconds,
		RefreshToken:     refreshToken,
		RefreshExpiresIn: m.refreshTTLSeconds,
		TokenType:        "bearer",
	}, nil
}

func (m *Manager) signedToken(userID int64, exp time.Time) (string, error) {
	claims := claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *Manager) Parse(tokenStr string) (*claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &claims{},
		func(t *jwt.Token) (interface{}, error) { return m.secret, nil })
	if err != nil {
		return nil, err
	}
	return token.Claims.(*claims), nil
}
