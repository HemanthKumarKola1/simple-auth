package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	cache "github.com/HemanthKumarKola1/simple-auth/internal/cache"
	"github.com/HemanthKumarKola1/simple-auth/internal/handlers"
	usecase "github.com/HemanthKumarKola1/simple-auth/internal/middleware"
	repo "github.com/HemanthKumarKola1/simple-auth/internal/repo"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {

	dbConn, err := sql.Open("postgres", "postgresql://root:secret@localhost:5432/users?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	_, err = dbConn.Exec("CREATE TABLE users (username TEXT PRIMARY KEY, password TEXT)")
	if err != nil {
		fmt.Println("Error creating table:", err)
	} else {
		fmt.Println("Table created successfully")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err = rdb.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}

	rtc := cache.NewRevokedTokensCache(rdb)

	userRepo := repo.NewRepository(dbConn)
	authUsecase := usecase.NewAuthUseCase(userRepo, rtc)
	server := handlers.NewAuthServer(authUsecase)

	r := mux.NewRouter()
	r.HandleFunc("/signup", server.Signup).Methods(http.MethodPost)
	r.HandleFunc("/login", server.Login).Methods(http.MethodPost)
	r.HandleFunc("/refresh", server.RefreshJwt).Methods(http.MethodGet)
	r.HandleFunc("/revoke", server.Revoke).Methods(http.MethodPost)
	log.Fatal(http.ListenAndServe(":8080", r))

	dbConn.Exec("DROP TABLE IF EXISTS users")
}
