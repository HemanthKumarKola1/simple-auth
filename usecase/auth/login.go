package usecase

import (
	"errors"

	"hemanth.kola/simple-auth/models"
	"hemanth.kola/simple-auth/usecase/utils"
)

func (a *authUsecase) Login(user models.User) (string, error) {
	for _, u := range users {
		if u.Username == user.Username && u.Password == user.Password {
			accessToken, err := utils.GenerateTokens(user.Username)
			if err != nil {
				return "", err
			}
			return accessToken, nil
		}
	}
	return "", errors.New("invalid credentials")
}
