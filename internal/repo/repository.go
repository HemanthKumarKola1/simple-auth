package repo

import (
	"context"
	"database/sql"
	"fmt"

	db "github.com/HemanthKumarKola1/simple-auth/internal/db/sqlc"
	"github.com/HemanthKumarKola1/simple-auth/internal/utils"
)

type Repository struct {
	*db.Queries
	db *sql.DB
}

func NewRepository(dbConn *sql.DB) *Repository {
	return &Repository{
		Queries: db.New(dbConn),
		db:      dbConn,
	}
}

func (r *Repository) CreateNewUser(ctx context.Context, args db.User) (db.User, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return db.User{}, err
	}
	defer tx.Rollback()

	user, err := r.GetUser(ctx, args.Username)
	if err != nil && err != sql.ErrNoRows {
		return db.User{}, err
	}
	if user.Username != "" {
		return db.User{}, fmt.Errorf(
			utils.ERROR_1,
		)
	}

	newUserArgs := db.CreateUserParams{
		Username: args.Username,
		Password: args.Password,
	}

	createdUser, err := r.CreateUser(ctx, newUserArgs)
	if err != nil {
		return db.User{}, err
	}

	if err := tx.Commit(); err != nil {
		return db.User{}, err
	}

	return createdUser, nil
}
