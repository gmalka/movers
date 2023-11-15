package rest

import (
	"net/http"

	"context"
	"log"

	"github.com/gmalka/movers/model"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type UserRequest struct{}

type Log struct {
	Err *log.Logger
	Inf *log.Logger
}

func (l Log) Error(str string) {
	l.Err.Println(str)
}

func (l Log) Info(str string) {
	l.Inf.Println(str)
}

type Handler struct {
	game  GameIterator
	users UserService
	tasks TaskService
	auth  AuthService

	log Log
}

func NewHandler(game GameIterator, users UserService, tasks TaskService, auth AuthService, log Log) Handler {
	return Handler{
		game:  game,
		users: users,
		tasks: tasks,
		auth:  auth,
		log:   log,
	}
}

func (h Handler) Init() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", h.MainMenu)
	r.Get("/login", h.LoginTemplate)
	r.Post("/login", h.Login)

	r.Get("/register", h.RegisterTemplate)
	r.Post("/register", h.Regsiter)

	r.Get("/tasks", h.CreateTasksTemplate)
	//r.Post("/tasks", h.CreateTasks)

	// r.Route("/{username}", func(r chi.Router) {
	// 	r.Use(h.checkAccess)

	// 	r.Get("/", h.UserMenu)
	// 	r.Get("/tasks", h.GetCompletedTasks)
	// 	r.Post("/start", h.IterateGame)
	// 	r.Post("/update", h.DeleteUser)
	// })

	return r
}

// <----------------INTERFACES---------------->

type GameIterator interface {
	IterateWork(ctx context.Context, customername string, workernames []string, task model.Task) error
}

type UserService interface {
	NewCustomer(ctx context.Context, customer model.CustomerInfo) error
	GetCustomer(ctx context.Context, name string) (model.CustomerInfo, error)
	DeleteCustomer(ctx context.Context, name string) error
	NewWorker(ctx context.Context, worker model.WorkerInfo) error
	GetWorkers(ctx context.Context) ([]model.WorkerInfo, error)
	GetWorker(ctx context.Context, name string) (model.WorkerInfo, error)
	DeleteWorker(ctx context.Context, name string) error
}

type TaskService interface {
	GenerateTasks(ctx context.Context, tocreate int) error
	GetFirstTask(ctx context.Context) (model.Task, error)
	GetTasks(ctx context.Context) ([]model.Task, error)
	GetWorkerTasks(ctx context.Context, name string) ([]model.Task, error)
	FinishTask(ctx context.Context, workers []string, task model.Task) error
}

type AuthService interface {
	Login(ctx context.Context, username, password string) (model.Tokens, error)
	Register(ctx context.Context, user model.User) error
	CheckAccessToken(token string) (model.UserInfo, error)
	DeleteUser(ctx context.Context, name string) error
	UpdateRefreshToken(token string) (string, error)
	UpdateAccessToken(token string) (string, error)
}
