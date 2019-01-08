package auth

import (
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

const (
	// defaultTTL is a default token time to live.
	defaultTTL = 10 * 24 * time.Hour

	// defaultIssuer is a token issuer.
	defaultIssuer = "nott"
)

// Tokener issues tokens for user authentication.
type Tokener interface {
	Issue(User) (Token, error)
	Parse(token string) (id int, err error)
}

// Token is used for user authentication.
type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at,omitempty"`
}

// JWTokener issues JWT.
type JWTokener struct {
	secret []byte
	issuer string
	ttl    time.Duration
}

// NewJWTokener creates new JWT tokener.
func NewJWTokener(secret string) *JWTokener {
	return &JWTokener{
		secret: []byte(secret),
		issuer: defaultIssuer,
		ttl:    defaultTTL,
	}
}

// Issue issues new token for the given user.
func (t *JWTokener) Issue(user User) (Token, error) {
	var exp int64
	if t.ttl > 0 {
		exp = time.Now().Add(t.ttl).Unix()
	}
	claims := jwt.StandardClaims{
		Issuer:    t.issuer,
		Subject:   strconv.Itoa(user.ID),
		ExpiresAt: exp,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	s, err := token.SignedString(t.secret)
	if err != nil {
		return Token{}, errors.Wrap(err, "sign token string")
	}

	return Token{AccessToken: s, ExpiresAt: exp}, nil
}

// Parse parses, validates and gets user ID from JWT claims.
func (t *JWTokener) Parse(accessToken string) (int, error) {
	// Parse and validate
	claims := jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(
		accessToken,
		&claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return t.secret, nil
		},
	)
	if err != nil {
		return 0, errors.Wrap(err, "parse token")
	}
	if !token.Valid {
		return 0, errors.New("invalid token")
	}

	// Parse claims and get user ID
	if claims.Subject == "" {
		return 0, errors.New("sub field is empty")
	}
	id, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return 0, errors.New("id is not number")
	}
	return id, nil
}
