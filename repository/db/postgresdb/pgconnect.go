package postgresdb

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type DbConfig struct {
	Host     string
	Port     string
	User     string
	Dbname   string
	Password string
	Sslmode  string
}

func NewPostgresConnect(config DbConfig) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", config.Host, config.Port, config.User, config.Password, config.Dbname, config.Sslmode))
	if err != nil {
		return nil, fmt.Errorf("can't connect to bd: %v", err)
	}

	return db, nil
}