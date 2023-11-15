package authservice

import (
	"context"
	"fmt"
	"time"

	"github.com/gmalka/movers/model"
)

const (
	AccessToken = iota
	RefreshToken
)

type authService struct {
	us UserStore
	pm PasswordManager
	tm TokenManager
}

func NewAuthService(us UserStore, pm PasswordManager, tm TokenManager) *authService {
	return &authService{
		us: us,
		pm: pm,
		tm: tm,
	}
}

func (a *authService) CheckAccessToken(token string) (model.UserInfo, error) {
	info, err := a.tm.ParseToken(token, AccessToken)
	if err != nil {
		return model.UserInfo{}, fmt.Errorf("cant check token: %v", err)
	}

	return info, nil
}

func (a *authService) UpdateRefreshToken(token string) (string, error) {
	info, err := a.tm.ParseToken(token, RefreshToken)
	if err != nil {
		return "", fmt.Errorf("cant check token: %v", err)
	}

	token, err = a.tm.CreateToken(info, RefreshToken)
	if err != nil {
		return "", fmt.Errorf("cant create token: %v", err)
	}
	return token, nil
}

func (a *authService) UpdateAccessToken(token string) (string, error) {
	info, err := a.tm.ParseToken(token, RefreshToken)
	if err != nil {
		return "", fmt.Errorf("cant check token: %v", err)
	}

	token, err = a.tm.CreateToken(info, AccessToken)
	if err != nil {
		return "", fmt.Errorf("cant create token: %v", err)
	}
	return token, nil
}

func (a *authService) Login(ctx context.Context, username, password string) (model.Tokens, error) {
	user, err := a.us.GetUser(ctx, username)
	if err != nil {
		return model.Tokens{}, fmt.Errorf("cant login: %v", err)
	}

	err = a.pm.CheckPassword(password, user.Password)
	if err != nil {
		return model.Tokens{}, fmt.Errorf("cant login: %v", err)
	}

	access, err := a.tm.CreateToken(model.UserInfo{
		Name: user.Name,
		Role: user.Role,
	}, AccessToken)
	if err != nil {
		return model.Tokens{}, fmt.Errorf("cant create access token: %v", err)
	}

	refresh, err := a.tm.CreateToken(model.UserInfo{
		Name: user.Name,
		Role: user.Role,
	}, RefreshToken)
	if err != nil {
		return model.Tokens{}, fmt.Errorf("cant create access token: %v", err)
	}

	return model.Tokens{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (a *authService) Register(ctx context.Context, user model.User) error {
	var err error

	if user.Role == "Customer" {
		err = a.us.CheckForCustomerRole(ctx)
		if err != nil {
			return fmt.Errorf("cant create new customer: %v", err)
		}
	}

	user.Password, err = a.pm.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("cant register user: %v", err)
	}

	err = a.us.CreateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("cant register user %s: %v", user.Name, err)
	}

	return nil
}

func (a *authService) DeleteUser(ctx context.Context, name string) error {
	return a.us.DeleteUser(ctx, name)
}

func (a *authService) GetAccessTTL() time.Duration {
	return a.tm.GetAccessTTL()
}

func (a *authService) GetRefreshTTL() time.Duration {
	return a.tm.GetRefreshTTL()
}

// <----------------INTERFACES---------------->

type UserStore interface {
	CreateUser(ctx context.Context, user model.User) error
	DeleteUser(ctx context.Context, name string) error
	GetUser(ctx context.Context, name string) (model.User, error)
	CheckForCustomerRole(ctx context.Context) error
}

type PasswordManager interface {
	HashPassword(password string) (string, error)
	CheckPassword(verifiable, wanted string) error
}

type TokenManager interface {
	CreateToken(userinfo model.UserInfo, kind int) (string, error)
	ParseToken(inputToken string, kind int) (model.UserInfo, error)

	GetAccessTTL() time.Duration
	GetRefreshTTL() time.Duration
}
