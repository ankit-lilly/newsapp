package auth

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ankit-lilly/newsapp/pkg/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

const (
	accessTokenCookieName = "access-token"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token has expired")
)

type JwtClaims struct {
	Username string `json:"userName"`
	Id       int64  `json:"id"`
	jwt.RegisteredClaims
}

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type JwtService struct {
	JWTSecret       string
	TokenExpiration time.Duration
	CookieDomain    string
	IsDevelopment   bool
	AppName         string
}

func NewJwtService() *JwtService {
	config := config.LoadConfig()
	return &JwtService{
		JWTSecret:       config.JwtSecret,
		TokenExpiration: 24 * time.Hour,
		CookieDomain:    config.CookieDomain,
		IsDevelopment:   config.IsDev,
		AppName:         config.AppName,
	}
}

func (s *JwtService) GenerateToken(user User) (string, time.Time, error) {
	now := time.Now()
	expirationTime := now.Add(s.TokenExpiration)

	claims := &JwtClaims{
		Username: user.Username,
		Id:       user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: expirationTime},
			IssuedAt:  &jwt.NumericDate{Time: now},
			NotBefore: &jwt.NumericDate{Time: now},
			Issuer:    s.AppName,
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.JWTSecret))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("sign token: %w", err)
	}

	return tokenString, expirationTime, nil
}

func (s *JwtService) ValidateToken(tokenString string) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Name {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(s.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(*JwtClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims.ExpiresAt.Before(time.Now()) {
		return nil, ErrTokenExpired
	}

	return claims, nil
}

// Cookie Management
func (s *JwtService) GenerateTokenAndSetCookie(user User, c echo.Context) error {
	token, exp, err := s.GenerateToken(user)
	if err != nil {
		return fmt.Errorf("generate token: %w", err)
	}

	s.setTokenCookie(token, exp, c)
	return nil
}

func (s *JwtService) setTokenCookie(token string, expiration time.Time, c echo.Context) {
	cookie := &http.Cookie{
		Name:     accessTokenCookieName,
		Value:    token,
		Expires:  expiration,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Domain:   s.CookieDomain,
		Secure:   !s.IsDevelopment,
	}

	c.SetCookie(cookie)
}

func (s *JwtService) Logout(c echo.Context) error {
	cookie := &http.Cookie{
		Name:     accessTokenCookieName,
		Value:    "",
		Expires:  time.Now().Add(-24 * time.Hour),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Domain:   s.CookieDomain,
		Secure:   !s.IsDevelopment,
	}

	c.SetCookie(cookie)
	return nil
}
