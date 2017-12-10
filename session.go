package winter

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"strings"
)

type Handler struct {
	token *string
	Store
}

type Session interface {
	New(payload Manifest, password ...string) error
	Authenticate() bool
}

type Manifest interface {
	Issuer() string
	Subject() string
	Expiry() int64
}

func key(handler *Handler, issuer string, subject string) *rsa.PrivateKey {
	uniqueId := strings.Join([]string{issuer, subject}, ":")
	return handler.Key[uniqueId]
	//handler.Key[issuer:subject]
}

func (handler *Handler) New(payload Manifest, password ...string) {
	issuer := payload.Issuer()
	subject := payload.Subject()
	claims := &jwt.StandardClaims{
		ExpiresAt: payload.Expiry(),
		Issuer:    issuer,
		Subject:   subject,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	privateKey := key(handler, issuer, subject)
	signedToken, _ := token.SignedString(*privateKey)

	handler.token = &signedToken
}
