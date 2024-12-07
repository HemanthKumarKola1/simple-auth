package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	db "hemanth.kola/simple-auth/internal/db/sqlc"
	usecase "hemanth.kola/simple-auth/internal/middleware"
	"hemanth.kola/simple-auth/internal/utils"
)

type authServer struct {
	authenticateUsecase usecase.Auth
}

type AuthServer interface {
	Signup(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	RefreshJwt(w http.ResponseWriter, r *http.Request)
	Revoke(w http.ResponseWriter, r *http.Request)
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
	token, err := utils.ExtractJWTFromHeader(r)
	if err != nil {
		http.Error(w, "error while extracting jwt from header"+err.Error(), http.StatusBadRequest)
		return
	}
	jwtToken, err := utils.ValidateJWT(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//TODO: check if the current token is revoked already.

	claims, err := utils.GetClaims(jwtToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	expiry := claims["exp"].(float64)

	if err := a.authenticateUsecase.RevokeJwt(token, expiry); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "token revoked successfully")
}