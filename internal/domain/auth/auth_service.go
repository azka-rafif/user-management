package auth

import (
	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/internal/domain/user"
	"github.com/evermos/boilerplate-go/shared/jwt"
)

type AuthService interface {
	Register(payload AuthPayload) (res JwtResponseFormat, err error)
	Login(payload LoginPayload) (res JwtResponseFormat, err error)
}

type AuthServiceImpl struct {
	Repo        AuthRepository
	Config      *configs.Config
	UserService user.UserService
}

func ProvideAuthServiceImpl(repo AuthRepository, conf *configs.Config, userService user.UserService) *AuthServiceImpl {
	return &AuthServiceImpl{Config: conf, Repo: repo, UserService: userService}
}

func (s *AuthServiceImpl) Register(payload AuthPayload) (res JwtResponseFormat, err error) {

	user, err := s.UserService.Create(user.UserPayload(payload))
	if err != nil {
		return
	}

	res, err = s.createToken(user)
	if err != nil {
		return
	}

	return
}

func (s *AuthServiceImpl) Login(payload LoginPayload) (res JwtResponseFormat, err error) {
	user, err := s.UserService.GetByUserName(payload.UserName)
	if err != nil {
		return
	}
	err = user.ValidatePassword(payload.Password)
	if err != nil {
		return
	}

	res, err = s.createToken(user)
	if err != nil {
		return
	}

	return
}

func (s *AuthServiceImpl) createToken(user user.User) (res JwtResponseFormat, err error) {
	jwt := jwt.NewJWT(s.Config.App.JWTSecret)
	token, err := jwt.GenerateJwt(user.UserId.String(), user.UserName, user.Role, user.CartId.String())
	if err != nil {
		return
	}
	res = JwtResponseFormat{AccessToken: token}
	return
}
