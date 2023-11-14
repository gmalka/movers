package userrepository

import (
	"context"
	"errors"
	"fmt"

	"github.com/gmalka/movers/model"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) CreateUser(ctx context.Context, user model.User) error {
	_, err := u.db.ExecContext(ctx, "INSERT INTO users VALUES($1,$2,$3)", user.Name, user.Password, user.Role)
	if err != nil {
		return fmt.Errorf("can't insert into users: %v", err)
	}

	return nil
}

func (u *UserRepository) GetUser(ctx context.Context, name string) (model.User, error) {
	row := u.db.QueryRowContext(ctx, "SELECT * FROM users WHERE name=$1", name)
	if row.Err() != nil {
		return model.User{}, fmt.Errorf("can't find user %s: %v", name, row.Err())
	}

	user := model.User{}
	err := row.Scan(&user.Name, &user.Password)
	if err != nil {
		return model.User{}, fmt.Errorf("can't scan user: %v", row.Err())
	}

	return user, nil
}

func (u *UserRepository) CheckForCustomerRole(ctx context.Context, ) error {
	var name string
	row := u.db.QueryRowContext(ctx, "SELECT name FROM users WHERE role=$1", "customer")

	err := row.Scan(&name)
	if err != nil {
		return errors.New("customer already exists")
	}

	return nil
}

func (u *UserRepository) DeleteUser(ctx context.Context, name string)  error {
	_, err := u.db.ExecContext(ctx, "DELETE FROM users WHERE name=$1", name)
	if err != nil {
		return fmt.Errorf("cant delete user %s: %v", name, err)
	}

	return nil
}