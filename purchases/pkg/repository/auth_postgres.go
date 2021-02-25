package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	templates "purchases/pkg"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{
		db: db,
	}
}

func (r *AuthPostgres) CreateUser(user templates.User) (int, error) {
	var id int

	query := fmt.Sprintf("INSERT INTO %s (name, password_hash) values ($1, $2) RETURNING id", usersTable)
	row := r.db.QueryRow(query, user.Name, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthPostgres) GetUser(username, password string) (templates.User, error) {
	var user templates.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE name=$1 AND password_hash=$2", usersTable)
	err := r.db.Get(&user, query, username, password)
	return user, err
}
