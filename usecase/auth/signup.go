package usecase

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"hemanth.kola/simple-auth/models"
)

var users = []models.User{
	{Username: "user1", Password: "password1"},
}

type authUsecase struct {
}

type Auth interface {
	SignUp(newUser *models.User) error
	Login(user models.User) (string, error)
	RefreshJwt(jwtToken string) (string, error)
	RevokeJwt(token string)
}

func (a *authUsecase) SignUp(newUser *models.User) error {
	// Check if username already exists
	for _, user := range users {
		if user.Username == newUser.Username {
			return errors.New("username already exists")
		}
	}

	// Hash the password (replace this with a strong hashing algorithm)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("error saving password")
	}

	newUser.Password = string(hashedPassword)

	users = append(users, *newUser)
	return nil
}
