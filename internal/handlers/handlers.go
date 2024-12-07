package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	db "github.com/HemanthKumarKola1/simple-auth/internal/db/sqlc"
	usecase "github.com/HemanthKumarKola1/simple-auth/internal/middleware"
	"github.com/HemanthKumarKola1/simple-auth/internal/utils"
	"github.com/golang-jwt/jwt"
)

type authServer struct {
	authenticateUsecase usecase.Auth
}

type AuthServer interface {
	Signup(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	RefreshJwt(w http.ResponseWriter, r *http.Request)
	Revoke(w http.ResponseWriter, r *http.Request)
	TestAuth(w http.ResponseWriter, r *http.Request)
}

func NewAuthServer(authUsecase usecase.Auth) AuthServer {
	return &authServer{authenticateUsecase: authUsecase}
}

func (a *authServer) Signup(w http.ResponseWriter, r *http.Request) {
	var newUser db.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "unable to decode request body", http.StatusBadRequest)
		return
	}
	err = a.authenticateUsecase.SignUp(newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(newUser)
}

func (a *authServer) Login(w http.ResponseWriter, r *http.Request) {
	var user db.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "unable to decode request body", http.StatusBadRequest)
		return
	}

	// Authenticate user
	accessToken, err := a.authenticateUsecase.Login(user)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"jwt": accessToken,
		})
		return
	}

	http.Error(w, err.Error(), http.StatusUnauthorized)
}

func (a *authServer) RefreshJwt(w http.ResponseWriter, r *http.Request) {

	jwtToken, err := utils.ExtractJWTFromHeader(r)
	if err != nil {
		http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
		return
	}

	accessToken, err := a.authenticateUsecase.RefreshJwt(jwtToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Send response with tokens
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"jwt": accessToken,
	})

}

func (a *authServer) Revoke(w http.ResponseWriter, r *http.Request) {
	jwtToken, err := utils.ExtractJWTFromHeader(r)
	if err != nil {
		http.Error(w, "error while extracting jwt from header"+err.Error(), http.StatusBadRequest)
		return
	}
	token, err := utils.ValidateJWT(jwtToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims := token.Claims.(jwt.MapClaims)

	expiry := claims["exp"].(float64)
	if expiry < float64(time.Now().Unix()) {
		http.Error(w, "token expired", http.StatusUnauthorized)
		return
	}

	if err := a.authenticateUsecase.RevokeJwt(jwtToken); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "token revoked successfully")
}

func (a *authServer) TestAuth(w http.ResponseWriter, r *http.Request) {
	jwtToken, err := utils.ExtractJWTFromHeader(r)
	if err != nil {
		http.Error(w, "error while extracting jwt from header"+err.Error(), http.StatusBadRequest)
		return
	}
	token, err := utils.ValidateJWT(jwtToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	claims := token.Claims.(jwt.MapClaims)

	expiry := claims["exp"].(float64)
	if expiry < float64(time.Now().Unix()) {
		http.Error(w, "token expired", http.StatusUnauthorized)
		return
	}

	if err := a.authenticateUsecase.IsRevoked(jwtToken); err != nil {
		if err.Error() == utils.ERROR_2 {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "authorized successfully")
}
