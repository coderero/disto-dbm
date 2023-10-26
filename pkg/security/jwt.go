package security

import (
	"crypto/rsa"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

func init() {
	private, err := os.ReadFile("./certs/private.key")
	if err != nil {
		panic(err)
	}
	public, err1 := os.ReadFile("./certs/public.pem")
	if err1 != nil {
		panic(err1)
	}
	key, err2 := jwt.ParseRSAPrivateKeyFromPEM(private)
	if err2 != nil {
		panic(err2)
	}
	pub, err3 := jwt.ParseRSAPublicKeyFromPEM(public)
	if err3 != nil {
		panic(err3)
	}
	privateKey, publicKey = key, pub

}

func GenerateToken(sub string, exp time.Time) string {
	expiresIn := jwt.NewNumericDate(exp)

	claims := jwt.RegisteredClaims{
		Subject:   sub,
		ExpiresAt: expiresIn,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		panic(err)
	}
	return signedToken
}

func VerifyToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
}
