package rest

import (
	"net/http"
	"strings"
	"time"

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
	game  WorkService
	users UserService
	tasks TaskService
	auth  AuthService

	log Log
}

func NewHandler(game WorkService, users UserService, tasks TaskService, auth AuthService, log Log) Handler {
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
	r.Post("/tasks", h.CreateTasks)

	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})

	r.Post("/refresh", h.refreshAccessToken)
	r.Post("/access", h.refreshRefreshToken)

	r.Route("/{username}", func(r chi.Router) {
		r.Use(h.checkAccess)

		r.Get("/", h.UserMenu)
		r.Get("/me", h.AboutMe)
		r.Get("/tasks", h.GetTasks)
		r.Get("/start", h.StartGame)
		r.Post("/iterate", h.ChoosewWorkers)
		r.Delete("/delete", h.DeleteUser)
		r.Get("/exit", h.Exit)
	})

	return r
}

func (h Handler) HandlerError(w http.ResponseWriter, err error, status int) {
	h.log.Error(err.Error())
	http.Error(w, http.StatusText(status), status)
}

func StepBack(url string) string {
	str := strings.Split(url, "/")
	l := len(str)

	l -= 1
	if l < 0 {
		l = 0
	}

	url = strings.Join(str[:l], "/")

	return url
}

// <----------------INTERFACES---------------->

type WorkService interface {
	MakeTaskForCustomer(ctx context.Context, customername string) error
}

type UserService interface {
	NewCustomer(ctx context.Context, customer model.CustomerInfo) error
	GetCustomer(ctx context.Context, name string) (model.CustomerInfo, error)
	DeleteCustomer(ctx context.Context, name string) error
	NewWorker(ctx context.Context, worker model.WorkerInfo) error
	GetWorkers(ctx context.Context) ([]model.WorkerInfo, error)
	GetWorker(ctx context.Context, name string) (model.WorkerInfo, error)

	GetChoosenWorkers(ctx context.Context) ([]model.WorkerInfo, error)
	RechooseWorkers(ctx context.Context, workers []string) error
}

type TaskService interface {
	GenerateTasks(ctx context.Context, tocreate int) error
	GetFirstTask(ctx context.Context) (model.Task, error)
	GetTasks(ctx context.Context, page int) ([]model.Task, error)
	GetWorkerTasks(ctx context.Context, name string, page int) ([]model.Task, error)
}

type AuthService interface {
	Login(ctx context.Context, username, password string) (model.Tokens, error)
	Register(ctx context.Context, user model.User) error
	CheckAccessToken(token string) (model.UserInfo, error)
	DeleteUser(ctx context.Context, name string) error
	UpdateRefreshToken(token string) (string, error)
	UpdateAccessToken(token string) (string, error)

	GetAccessTTL() time.Duration
}
