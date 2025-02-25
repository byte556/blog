package service

import (
	"Blog/models"
	"Blog/pkg/repository"
	"golang.org/x/crypto/bcrypt"
)

func (serv *AuthService) SignIn(username, password string) (int, error) {
	user, err := serv.GetUserByName(username)
	if err != nil {
		return 0, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return 0, err
	}
	return user.Id, nil
}

func (serv *AuthService) SignUp(username, password string) error {
	if username == "" {
		return models.ErrUsernameIsEmpty
	}
	if password == "" {
		return models.ErrPasswordIsEmpty
	}

	_, err := serv.GetUserByName(username)
	if err != models.ErrUserNotFound {
		return models.ErrUserExists
	}

	psswrd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}
	newUser := &models.User{Name: username, Password: string(psswrd)}
	return serv.AddUser(newUser)

}

func (serv *AuthService) GetUser(id int) (models.User, error) {
	return serv.GetUserById(id)
}

type AuthService struct {
	repository.Authorization
}

func NewAuthService(authorization repository.Authorization) *AuthService {
	return &AuthService{authorization}
}
