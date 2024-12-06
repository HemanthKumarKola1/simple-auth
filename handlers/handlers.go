package handlers

import (
	"encoding/json"
	"net/http"

	"hemanth.kola/simple-auth/models"
	usecase "hemanth.kola/simple-auth/usecase/auth"
	"hemanth.kola/simple-auth/utils"
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

func NewAuthServer() AuthServer {
	return &authServer{}
}

func (a *authServer) Signup(w http.ResponseWriter, r *http.Request) {
	var newUser models.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "unable to decode request body", http.StatusBadRequest)
		return
	}
	err = a.authenticateUsecase.SignUp(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(newUser)
}

func (a *authServer) Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
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

	jwtToken, err := utils.ExtractJWT(r)
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
	jwtToken, err := utils.ExtractJWT(r)
	if err != nil {
		http.Error(w, "error while extracting jwt from header"+err.Error(), http.StatusBadRequest)
	}
	a.authenticateUsecase.RevokeJwt(jwtToken)
}
