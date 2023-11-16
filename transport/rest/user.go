package rest

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/gmalka/movers/model"
	"github.com/go-chi/chi"
)

func (h Handler) UserMenu(w http.ResponseWriter, r *http.Request) {
	fp := path.Join(h.PathToTemplates()+"/templates/user", "menu.html")

	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		h.HandlerError(w, err, http.StatusInternalServerError)
		return
	}

	u, ok := r.Context().Value(UserRequest{}).(model.UserInfo)
	if !ok {
		h.HandlerError(w, err, http.StatusInternalServerError)
		return
	}

	switch u.Role {
	case "Customer":
		err = tmpl.ExecuteTemplate(w, "menu", true)
	case "Worker":
		err = tmpl.ExecuteTemplate(w, "menu", nil)
	default:
		err = fmt.Errorf("unknown role: %v", u.Role)
	}
	if err != nil {
		h.HandlerError(w, err, http.StatusInternalServerError)
		return
	}
}

func (h Handler) AboutMe(w http.ResponseWriter, r *http.Request) {
	fp := path.Join(h.PathToTemplates()+"/templates/user", "aboutme.html")

	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		h.HandlerError(w, err, http.StatusInternalServerError)
		return
	}

	u, ok := r.Context().Value(UserRequest{}).(model.UserInfo)
	if !ok {
		h.HandlerError(w, err, http.StatusInternalServerError)
		return
	}
	switch u.Role {
	case "Customer":
		data := struct {
			Customer bool
			Name     string
			Money    int
			Workers  []model.WorkerInfo
		}{
			Name:     u.Name,
			Customer: true,
		}
		customer, err := h.users.GetCustomer(r.Context(), u.Name)
		if err != nil {
			h.HandlerError(w, err, http.StatusBadRequest)
			return
		}

		data.Money = customer.Money
		data.Workers, err = h.users.GetChoosenWorkers(r.Context())
		if err != nil {
			h.HandlerError(w, err, http.StatusInternalServerError)
			return
		}
		err = tmpl.ExecuteTemplate(w, "aboutme", data)
		if err != nil {
			h.HandlerError(w, err, http.StatusInternalServerError)
			return
		}
	case "Worker":
		worker, err := h.users.GetWorker(r.Context(), u.Name)
		if err != nil {
			h.HandlerError(w, err, http.StatusBadRequest)
			return
		}

		data := struct {
			Customer    bool
			Name        string
			Fatigue     int
			Salary      int
			CarryWeight int
			Drunk       bool
		}{
			Customer:    false,
			Name:        worker.Name,
			Fatigue:     worker.Fatigue,
			Salary:      worker.Salary,
			CarryWeight: worker.CarryWeight,
		}
		if worker.Drunk == 2 {
			data.Drunk = true
		}

		err = tmpl.ExecuteTemplate(w, "aboutme", data)
		if err != nil {
			h.HandlerError(w, err, http.StatusInternalServerError)
			return
		}
	default:
		if err != nil {
			h.HandlerError(w, err, http.StatusBadRequest)
			return
		}
	}
}

func (h Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	u, ok := r.Context().Value(UserRequest{}).(model.UserInfo)
	if !ok {
		h.HandlerError(w, errors.New("cant take UserInfo from context"), http.StatusInternalServerError)
		return
	}

	fp := path.Join(h.PathToTemplates()+"/templates/user", "tasks.html")

	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		h.HandlerError(w, err, http.StatusInternalServerError)
		return
	}

	var data []model.Task

	p := r.URL.Query().Get("page")
	page := 1
	if p != "" {
		page, err = strconv.Atoi(p)
		if err != nil {
			h.HandlerError(w, err, http.StatusBadRequest)
			return
		}
	}
	if page <= 0 {
		page = 1
	}

	if u.Role == "Customer" {
		data, err = h.tasks.GetTasks(r.Context(), page)
		if err != nil {
			h.HandlerError(w, err, http.StatusBadRequest)
			return
		}
	} else {
		data, err = h.tasks.GetWorkerTasks(r.Context(), u.Name, page)
		if err != nil {
			h.HandlerError(w, err, http.StatusBadRequest)
			return
		}
	}

	if err := tmpl.ExecuteTemplate(w, "tasks", struct {
		Page  int
		Tasks []model.Task
	}{
		Page:  page,
		Tasks: data,
	}); err != nil {
		h.HandlerError(w, err, http.StatusInternalServerError)
		return
	}
}

func (h Handler) StartGame(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "username")

	err := r.ParseForm()
	if err != nil {
		h.HandlerError(w, err, http.StatusInternalServerError)
		return
	}

	workers, err := h.users.GetWorkers(r.Context())
	if err != nil {
		h.HandlerError(w, err, http.StatusInternalServerError)
		return
	}

	customer, err := h.users.GetCustomer(r.Context(), name)
	if err != nil {
		h.HandlerError(w, err, http.StatusBadRequest)
		return
	}

	fp := path.Join(h.PathToTemplates()+"/templates/user", "workerchoose.html")

	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		h.HandlerError(w, err, http.StatusInternalServerError)
		return
	}

	task, err := h.tasks.GetFirstTask(r.Context())
	if err != nil {
		h.HandlerError(w, err, http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "workerchoose", struct {
		Lost     bool
		ItemName string
		Weight   int
		Message  string
		Money    int
		Workers  []model.WorkerInfo
	}{
		Lost:     customer.Lost,
		ItemName: task.ItemName,
		Weight:   task.Weight,
		Message:  "",
		Money:    customer.Money,
		Workers:  workers,
	}); err != nil {
		h.log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h Handler) ChoosewWorkers(w http.ResponseWriter, r *http.Request) {
	message := "Success!"

	err := r.ParseForm()
	if err != nil {
		h.HandlerError(w, err, http.StatusInternalServerError)
		return
	}

	v, ok := r.Form["selectedWorkers"]
	if !ok {
		h.log.Info("No worker choosen")
		v = []string{}
	}

	err = h.users.RechooseWorkers(r.Context(), v)
	if err != nil {
		h.HandlerError(w, err, http.StatusBadRequest)
		return
	}

	name := chi.URLParam(r, "username")

	if len(v) != 0 {
		err = h.game.MakeTaskForCustomer(r.Context(), name)
		if err != nil {
			message = fmt.Sprintf("Game ended, %v", err.Error())
			h.log.Error(err.Error())
		}
	} else {
		message = "No one worker was choosen"
	}

	fp := path.Join(h.PathToTemplates()+"/templates/user", "workerchoose.html")

	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		h.HandlerError(w, err, http.StatusInternalServerError)
		return
	}

	customer, err := h.users.GetCustomer(r.Context(), name)
	if err != nil {
		h.HandlerError(w, err, http.StatusBadRequest)
		return
	}

	workers, err := h.users.GetWorkers(r.Context())
	if err != nil {
		h.HandlerError(w, err, http.StatusInternalServerError)
		return
	}

	task, err := h.tasks.GetFirstTask(r.Context())
	if err != nil {
		message = "no tasks in the roll"
		h.log.Error(err.Error())
	}

	if err := tmpl.ExecuteTemplate(w, "workerchoose", struct {
		Lost     bool
		ItemName string
		Weight   int
		Message  string
		Money    int
		Workers  []model.WorkerInfo
	}{
		Lost:     customer.Lost,
		ItemName: task.ItemName,
		Weight:   task.Weight,
		Message:  message,
		Money:    customer.Money,
		Workers:  workers,
	}); err != nil {
		h.HandlerError(w, err, http.StatusInternalServerError)
		return
	}
}

func (h Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	u, ok := r.Context().Value(UserRequest{}).(model.UserInfo)
	if !ok {
		h.HandlerError(w, errors.New("cant take UserInfo from context"), http.StatusInternalServerError)
		return
	}

	if u.Role != "Customer" {
		h.HandlerError(w, errors.New("try to delete worker"), http.StatusBadRequest)
		return
	}

	err := h.auth.DeleteUser(r.Context(), u.Name)
	if err != nil {
		h.HandlerError(w, err, http.StatusBadRequest)
		return
	}

	err = h.users.DeleteCustomer(r.Context(), u.Name)
	if err != nil {
		h.HandlerError(w, err, http.StatusBadRequest)
		return
	}

	cookie := CreateBearerTokenCookie("", "access_token", "/"+u.Name, time.Now().Add(-1*time.Hour))
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h Handler) Exit(w http.ResponseWriter, r *http.Request) {
	u, ok := r.Context().Value(UserRequest{}).(model.UserInfo)
	if !ok {
		h.HandlerError(w, errors.New("cant take UserInfo from context"), http.StatusInternalServerError)
		return
	}

	cookie := CreateBearerTokenCookie("", "access_token", "/"+u.Name, time.Now().Add(-1*time.Hour))
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
