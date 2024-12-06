package usecase

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"hemanth.kola/simple-auth/usecase/utils"
)

var revokedJwts = map[string]bool{
	"klasjfkl": false,
}

func (a *authUsecase) RefreshJwt(jwtToken string) (string, error) {
	token, err := utils.ValidateJWT(jwtToken)
	if err != nil {
		return "", errors.New("invalid jwt provided")
	}

	if _, ok := revokedJwts[jwtToken]; ok {
		return "", errors.New("invalid jwt provided")
	}

	claims := token.Claims.(jwt.MapClaims)
	uname := claims["username"].(string)
	expiry := claims["exp"].(int64)
	if expiry < int64(time.Now().Unix()) {
		return "", errors.New("token expired")
	}

	accessToken, err := utils.GenerateTokens(uname)
	if err != nil {
		return "", errors.New("unable to create token, try again")
	}

	return accessToken, nil
}
