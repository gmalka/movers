package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gmalka/movers/repository/db/postgresdb"
	"github.com/gmalka/movers/repository/db/postgresdb/customerrepository"
	"github.com/gmalka/movers/repository/db/postgresdb/donetasksrepository"
	"github.com/gmalka/movers/repository/db/postgresdb/itemrepository"
	"github.com/gmalka/movers/repository/db/postgresdb/taskrepository"
	"github.com/gmalka/movers/repository/db/postgresdb/userrepository"
	"github.com/gmalka/movers/repository/db/postgresdb/workerrepository"
	"github.com/gmalka/movers/service/auth/authservice"
	"github.com/gmalka/movers/service/auth/passwordservice"
	tokenmanager "github.com/gmalka/movers/service/auth/tokenservice"
	"github.com/gmalka/movers/service/taskservice"
	"github.com/gmalka/movers/service/userinfoservice"
	"github.com/gmalka/movers/service/workservice"
	"github.com/gmalka/movers/transport/rest"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	conn, err := postgresdb.NewPostgresConnect(
		postgresdb.Host("db"),
		postgresdb.Port("5432"),
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	customers := customerrepository.NewCustomerRepository(conn)
	donetask := donetasksrepository.NewDoneTasksRepository(conn)
	items := itemrepository.NewItemService(conn)
	tasks := taskrepository.NewTaskRepository(conn)
	users := userrepository.NewUserRepository(conn)
	worker := workerrepository.NewWorkerRepository(conn)

	passwords := passwordservice.NewPasswordManager()
	tokens := tokenmanager.NewAuthService(os.Getenv("ACCESS_SECRET"), os.Getenv("REFRESH_SECRET"))
	auth := authservice.NewAuthService(users, passwords, tokens)
	taskservice := taskservice.NewTaskService(items, tasks, donetask)
	userinfo := userinfoservice.NewuserInfoService(customers, worker)
	workservice := workservice.NewWorkService(customers, worker, taskservice)

	h := rest.NewHandler(workservice, userinfo, taskservice, auth, rest.Log{
		Err: log.Default(),
		Inf: log.Default(),
	})

	server := http.Server{
		Addr:    ":8080",
		Handler: h.Init(),
	}

	fmt.Println(server.ListenAndServe())
}
