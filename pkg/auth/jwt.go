package auth

import (
	"fmt"
	"net/http"
	"time"

	repository "github.com/ankibahuguna/newsapp/internal/auth/respository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

const (
	accessTokenCookieName = "access-token"
	jwtSecretKey          = "mySuperSecretComplecatedJWTKeyThat's1500xTimesHardToCrackThanAKeyLike-jwt-secret"
)

var TokenLookUpKey = fmt.Sprintf("cookie:%s", accessTokenCookieName)

func GetJWTSecret() string {
	return jwtSecretKey
}

type JwtClaims struct {
	Name string `json:"name"`
	Id   int64  `json:"id"`
	jwt.RegisteredClaims
}

func GenerateTokensAndSetCookies(user *repository.User, expTime time.Duration, c echo.Context) error {
	expirationTime := time.Now().Add(expTime)
	accessToken, exp, err := GenerateAccessToken(user, expirationTime, []byte(jwtSecretKey))
	if err != nil {
		return err
	}

	setTokenCookie(accessTokenCookieName, accessToken, exp, c)

	return nil
}

func GenerateAccessToken(user *repository.User, expirationTime time.Time, secret []byte) (string, time.Time, error) {
	claims := &JwtClaims{
		Name: user.Name,
		Id:   user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{expirationTime},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string.
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", time.Now(), err
	}

	return tokenString, expirationTime, nil
}

// Here we are creating a new cookie, which will store the valid JWT token.
func setTokenCookie(name, token string, expiration time.Time, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = token
	cookie.Expires = expiration
	cookie.SameSite = http.SameSiteStrictMode
	cookie.Path = "/"
	cookie.HttpOnly = true

	c.SetCookie(cookie)
}

func JWTErrorChecker(err error, c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, c.Echo().Reverse("userSignInForm"))
}
