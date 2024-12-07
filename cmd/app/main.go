package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/HemanthKumarKola1/simple-auth/internal/cache"
	"github.com/HemanthKumarKola1/simple-auth/internal/handlers"
	middleware "github.com/HemanthKumarKola1/simple-auth/internal/middleware"
	"github.com/HemanthKumarKola1/simple-auth/internal/repo"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {

	time.Sleep(5 * time.Second) // to make the dependencies get up meanwhile

	pgsqlUrl := "postgresql://root:secret@db:5432/users?sslmode=disable"
	dbConn, err := sql.Open("postgres", pgsqlUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err = rdb.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}

	_, err = dbConn.Exec("CREATE TABLE users (username TEXT PRIMARY KEY, password TEXT)")
	if err != nil {
		fmt.Println("Error creating table:", err)
	} else {
		fmt.Println("Table created successfully")
	}

	rtc := cache.NewRevokedTokensCache(rdb)

	userRepo := repo.NewRepository(dbConn)
	authUsecase := middleware.NewAuthUseCase(userRepo, rtc)
	server := handlers.NewAuthServer(authUsecase)

	r := mux.NewRouter()
	r.HandleFunc("/signup", server.Signup).Methods(http.MethodPost)
	r.HandleFunc("/login", server.Login).Methods(http.MethodPost)
	r.HandleFunc("/refresh", server.RefreshJwt).Methods(http.MethodGet)
	r.HandleFunc("/revoke", server.Revoke).Methods(http.MethodPost)
	r.HandleFunc("/test", server.TestAuth).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(":8080", r))
}
