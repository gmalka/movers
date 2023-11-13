package authservice

import (
	"fmt"

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

func NewAuthService(us UserStore, pm PasswordManager, tm TokenManager) authService {
	return authService{
		us: us,
		pm: pm,
		tm: tm,
	}
}

func (a authService) Login(username, password string) (model.Tokens, error) {
	user, err := a.us.GetUser(username)
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
		AccessToken: access,
		RefreshToken: refresh,
	}, nil
}

func (a authService) Register(username, password, role string) error {
	var err error

	password, err = a.pm.HashPassword(password)
	if err != nil {
		return fmt.Errorf("cant register user: %v", err)
	}

	
}

type UserStore interface {
	CreateUser(user model.User) error
	GetUser(name string) (model.User, error)
}

type PasswordManager interface {
	HashPassword(password string) (string, error)
	CheckPassword(verifiable, wanted string) error
}

type TokenManager interface {
	CreateToken(userinfo model.UserInfo, kind int) (string, error)
	ParseToken(inputToken string, kind int) (model.UserInfo, error)
}
