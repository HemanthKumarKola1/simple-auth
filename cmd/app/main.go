package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"hemanth.kola/simple-auth/handlers"
)

func main() {

	server := handlers.NewAuthServer()
	r := mux.NewRouter()
	r.HandleFunc("/singup", server.Signup).Methods(http.MethodPost)
	r.HandleFunc("/login", server.Login).Methods(http.MethodPost)
	r.HandleFunc("/refresh", server.RefreshJwt).Methods(http.MethodGet)
	r.HandleFunc("/revoke", server.Revoke).Methods(http.MethodPost)
	log.Fatal(http.ListenAndServe(":8080", r))
}
