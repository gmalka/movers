package postgresdb

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	defaultHost     = "localhost"
	defaultPort     = "5432"
	defaultUser     = "postgres"
	defaultDbname   = "postgres"
	defaultPassword = "postgres"
	defaultSslmode  = "disable"
)

type DbConfig struct {
	Host     string
	Port     string
	User     string
	Dbname   string
	Password string
	Sslmode  string
}

type Option func(*DbConfig)

func Host(host string) Option {
	return func(dc *DbConfig) {
		dc.Host = host
	}
}

func Port(port string) Option {
	return func(dc *DbConfig) {
		dc.Port = port
	}
}

func User(user string) Option {
	return func(dc *DbConfig) {
		dc.User = user
	}
}

func Dbname(dbname string) Option {
	return func(dc *DbConfig) {
		dc.Dbname = dbname
	}
}

func Password(password string) Option {
	return func(dc *DbConfig) {
		dc.Password = password
	}
}

func Sslmode(sslmode string) Option {
	return func(dc *DbConfig) {
		dc.Sslmode = sslmode
	}
}

func NewPostgresConnect(opts ...Option) (*sqlx.DB, error) {
	config := DbConfig{
		Host:     defaultHost,
		Port:     defaultPort,
		User:     defaultUser,
		Dbname:   defaultDbname,
		Password: defaultPassword,
		Sslmode:  defaultSslmode,
	}

	for _, v := range opts {
		v(&config)
	}

	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", config.Host, config.Port, config.User, config.Password, config.Dbname, config.Sslmode))
	if err != nil {
		return nil, fmt.Errorf("can't connect to bd: %v", err)
	}

	return db, nil
}
