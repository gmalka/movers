package tokenmanager

import (
	"errors"
	"fmt"
	"time"

	"github.com/gmalka/movers/model"
	"github.com/golang-jwt/jwt/v4"
)

const (
	ACCESS_TOKEN_TTL  = 15
	REFRESH_TOKEN_TTL = 60

	AccessToken = iota
	RefreshToken
)

type UserClaims struct {
	Name string `json:"name"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type authService struct {
	accessSecret  []byte
	refreshSecret []byte
}

// Прописать там, где будет использоваться
//
// type TokenManager interface {
// 	CreateToken(userinfo UserInfo, ttl time.Duration, kind int) (string, error)
// 	ParseToken(inputToken string, kind int) (UserInfo, error)
// }

func NewAuthService(accessSecret, refreshSecret string) *authService {
	return &authService{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
	}
}

func (u *authService) ParseToken(inputToken string, kind int) (model.UserInfo, error) {
	token, err := jwt.Parse(inputToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}

		var secret []byte
		switch kind {
		case AccessToken:
			secret = u.accessSecret
		case RefreshToken:
			secret = u.refreshSecret
		default:
			return "", fmt.Errorf("unknown secret kind %d", kind)
		}

		return secret, nil
	})

	if err != nil {
		return model.UserInfo{}, fmt.Errorf("can't parse token: %v", err)
	}

	if !token.Valid {
		return model.UserInfo{}, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return model.UserInfo{}, fmt.Errorf("can't get user claims from token")
	}

	return model.UserInfo{
		Name: claims["name"].(string),
		Role: claims["role"].(string),
	}, nil
}

func (u *authService) CreateToken(userinfo model.UserInfo, kind int) (string, error) {
	var ttl time.Duration

	switch kind {
	case AccessToken:
		ttl = ACCESS_TOKEN_TTL
	case RefreshToken:
		ttl = REFRESH_TOKEN_TTL
	default:
		return "", fmt.Errorf("unknown secret kind %d", kind)
	}

	claims := UserClaims{
		Name:             userinfo.Name,
		Role:             userinfo.Role,
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl * time.Minute))},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	var secret []byte
	switch kind {
	case AccessToken:
		secret = u.accessSecret
	case RefreshToken:
		secret = u.refreshSecret
	default:
		return "", fmt.Errorf("unknown secret kind %d", kind)
	}

	return token.SignedString(secret)
}
