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

type Claims struct {
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

	access, err := m.signedToken(userID, now.Add(time.Duration(m.accessTTLSeconds)*time.Second))
	if err != nil {
		return nil, err
	}
	refresh, err := m.signedToken(userID, now.Add(time.Duration(m.refreshTTLSeconds)*time.Second))
	if err != nil {
		return nil, err
	}

	return &Tokens{
		AccessToken:      access,
		ExpiresIn:        m.accessTTLSeconds,
		RefreshToken:     refresh,
		RefreshExpiresIn: m.refreshTTLSeconds,
		TokenType:        "bearer",
	}, nil
}

func (m *Manager) Parse(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{},
		func(t *jwt.Token) (interface{}, error) { return m.secret, nil })
	if err != nil {
		return nil, err
	}
	return token.Claims.(*Claims), nil
}

func (m *Manager) signedToken(userID int64, exp time.Time) (string, error) {
	cls := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, cls)
	return tkn.SignedString(m.secret)
}
