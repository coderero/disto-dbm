package security

import (
	"crypto/rsa"
	"net/http"
	"os"
	"time"

	"coderero.dev/projects/go/gin/hello/cache"
	"coderero.dev/projects/go/gin/hello/types"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// The `var` block is declaring two variables `privateKey` and `publicKey` of type `*rsa.PrivateKey`
// and `*rsa.PublicKey` respectively. These variables will be used to store the private and public keys
// for generating and verifying JWT tokens.
var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

// The `var` block is declaring two variables `AcessTokenExpireTime` and `RefreshTokenExpireTime` and
// assigning them values using the `time.Now().Add()` function. These variables will be used to store
// the expiration time for the access and refresh tokens.
var (
	AcessTokenExpireTime   = time.Now().Add(15 * time.Minute)
	RefreshTokenExpireTIme = time.Now().Add(7 * 24 * time.Hour)
)

func init() {

	// The code is reading the contents of two files, `private.key` and `public.pem`, located in the
	// `./certs` directory. It uses the `os.ReadFile()` function to read the files and assigns the contents
	// to the variables `private` and `public` respectively.
	private, err := os.ReadFile("./certs/private.key")
	if err != nil {
		panic(err)
	}
	public, err1 := os.ReadFile("./certs/public.pem")
	if err1 != nil {
		panic(err1)
	}

	// The code block is parsing the private and public keys from the PEM-encoded files `private.key` and
	// `public.pem` respectively.
	key, err2 := jwt.ParseRSAPrivateKeyFromPEM(private)
	if err2 != nil {
		panic(err2)
	}
	pub, err3 := jwt.ParseRSAPublicKeyFromPEM(public)
	if err3 != nil {
		panic(err3)
	}

	// The code block is assigning the parsed private and public keys to the variables `privateKey` and
	privateKey, publicKey = key, pub

}

// The function generates a JWT token with a specified subject and expiration time using the RS256
// signing method.
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

// The function `VerifyToken` parses a JWT token using a provided public key.
func VerifyToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
}

// The function `IsTokenExpired` checks if a given JWT token is expired or not.
func IsTokenExpired(token string) bool {
	jwtToken, err := VerifyToken(token)
	if err != nil {
		return true
	}
	if jwtToken.Valid {
		return false
	}
	return true
}

func IsTokenValid(token string) bool {
	jwtToken, err := VerifyToken(token)
	if err != nil {
		return false
	}
	if jwtToken.Valid {
		return true
	}
	return false
}

// The function checks if a token has been revoked and returns a boolean value indicating whether the
// token is revoked or not.
func TokenRevoked(accessToken string, refreshToken string, c *gin.Context, refresh bool) bool {
	if refresh {
		if cache.IsTokenRevoked(refreshToken) || cache.IsTokenRevoked(accessToken) {
			c.JSON(http.StatusUnauthorized, types.Response{
				Status: types.Status{
					Code: http.StatusUnauthorized,
					Msg:  "Any of the Token's Have been Revoked",
				},
			})
			return true
		}
	}

	if cache.IsTokenRevoked(accessToken) {
		c.JSON(http.StatusUnauthorized, types.Response{
			Status: types.Status{
				Code: http.StatusUnauthorized,
				Msg:  "Access Token has been Revoked",
			},
		})
		return true
	}
	return false
}
