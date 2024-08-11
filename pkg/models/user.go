package models

import (
	"context"

	"github.com/manosriram/outagealert.io/sqlc/db"
)

type UserModel struct {
	Db *db.Queries
}

func (u *UserModel) Create(user *db.User) error {
	_, err := u.Db.CreateUser(context.TODO(), db.CreateUserParams{
		Email:    user.Email,
		Password: user.Password,
	})
	return err

}

func (u *UserModel) All() ([]db.User, error) {
	users, _ := u.Db.AllUsers(context.TODO())
	return users, nil

}
