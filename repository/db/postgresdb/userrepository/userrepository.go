package userrepository

import (
	"fmt"

	"github.com/gmalka/movers/model"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return UserRepository{
		db: db,
	}
}

func (u *UserRepository) CreateUser(user model.User) error {
	_, err := u.db.Exec("INSERT INTO users VALUES($1,$2,$3)", user.Name, user.Password, user.Role)
	if err != nil {
		return fmt.Errorf("can't insert into users: %v", err)
	}
	
	return nil
}

func (u *UserRepository) CheckUser(name string) (model.User, error) {
	row := u.db.QueryRow("SELECT * FROM users WHERE name=$1", name)
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