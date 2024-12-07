package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	cache "hemanth.kola/simple-auth/internal/cache"
	db "hemanth.kola/simple-auth/internal/db/sqlc"
	repo "hemanth.kola/simple-auth/internal/repo"
	"hemanth.kola/simple-auth/internal/utils"
)

type Auth interface {
	SignUp(newUser db.User) error
	Login(user db.User) (string, error)
	RefreshJwt(jwtToken string) (string, error)
	RevokeJwt(token string, ttl float64) error
}

type authUsecase struct {
	repo              *repo.Repository
	revokeTokensCache *cache.RevokedTokensCache
}

func NewAuthUseCase(repo *repo.Repository, cache *cache.RevokedTokensCache) Auth {
	return &authUsecase{repo: repo, revokeTokensCache: cache}
}

func (a *authUsecase) SignUp(newUser db.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("error saving password")
	}

	newUser.Password = string(hashedPassword)

	if _, err := a.repo.CreateNewUser(context.Background(), newUser); err != nil {
		return err
	}
	return nil
}

func (a *authUsecase) Login(user db.User) (string, error) {

	dbUser, err := a.repo.GetUser(context.Background(), user.Username)
	if err != nil {
		return "", err
	}

	if err := comparePasswords(dbUser.Password, user.Password); err != nil {
		return "", errors.New("invalid credentials")
	}

	accessToken, err := utils.GenerateToken(user.Username)
	if err != nil {
		return "", err
	}
	return accessToken, nil

}

func comparePasswords(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}

func (a *authUsecase) RefreshJwt(jwtToken string) (string, error) {
	token, err := utils.ValidateJWT(jwtToken)
	if err != nil {
		return "", errors.New("invalid jwt provided")
	}

	val, err := a.revokeTokensCache.Get(jwtToken)
	if err == nil {
		return "", errors.New("invalid jwt provided")
	} else if err != redis.Nil {
		return "", errors.New("redis is unavailable")
	}

	fmt.Println(val)
	claims := token.Claims.(jwt.MapClaims)

	expiry := claims["exp"].(float64)
	if expiry < float64(time.Now().Unix()) {
		return "", errors.New("token expired")
	}

	uname := claims["username"].(string)
	accessToken, err := utils.GenerateToken(uname)
	if err != nil {
		return "", fmt.Errorf("unable to create token, err, %v", err.Error())
	}

	return accessToken, nil
}

func (a *authUsecase) RevokeJwt(token string, ttl float64) error {
	return a.revokeTokensCache.SetWithTTL(token, "revoked", 24*time.Hour)
}
