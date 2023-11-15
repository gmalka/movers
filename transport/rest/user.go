package rest

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
	"time"

	"github.com/gmalka/movers/model"
	"github.com/go-chi/chi"
)

func (h Handler) UserMenu(w http.ResponseWriter, r *http.Request) {
	fp := path.Join(h.PathToTemplates()+"/templates/user", "menu.html")

	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		h.log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	u, ok := r.Context().Value(UserRequest{}).(model.UserInfo)
	if !ok {
		h.log.Error("Cant take UserInfo from context")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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
		h.log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h Handler) AboutMe(w http.ResponseWriter, r *http.Request) {
	fp := path.Join(h.PathToTemplates()+"/templates/user", "aboutme.html")

	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		h.log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	u, ok := r.Context().Value(UserRequest{}).(model.UserInfo)
	if !ok {
		h.log.Error("Cant take UserInfo from context")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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
			h.log.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		data.Money = customer.Money
		data.Workers, err = h.users.GetChoosenWorkers(r.Context())
		if err != nil {
			h.log.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		err = tmpl.ExecuteTemplate(w, "aboutme", data)
		if err != nil {
			h.log.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	case "Worker":
		worker, err := h.users.GetWorker(r.Context(), u.Name)
		if err != nil {
			h.log.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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
			h.log.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	default:
		if err != nil {
			h.log.Error(fmt.Errorf("unknown role: %v", u.Role).Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

func (h Handler) GetCompletedTasks(w http.ResponseWriter, r *http.Request) {

}

func (h Handler) IterateGameGet(w http.ResponseWriter, r *http.Request) {
	r.Method = http.MethodPost
	http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
}

func (h Handler) IterateGame(w http.ResponseWriter, r *http.Request) {
}

func (h Handler) DeleteUserGet(w http.ResponseWriter, r *http.Request) {
	r.Method = http.MethodDelete
	http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
}

func (h Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "username")

	workers, err := h.users.GetChoosenWorkers(r.Context())
	if err != nil {
		h.log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	task, err := h.tasks.GetFirstTask(r.Context())
	if err != nil {
		h.log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	h.game.CalculateWork(r.Context(), name, workers, task)
}

func (h Handler) Exit(w http.ResponseWriter, r *http.Request) {
	u, ok := r.Context().Value(UserRequest{}).(model.UserInfo)
	if !ok {
		h.log.Error("Cant take UserInfo from context")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	cookie := CreateBearerTokenCookie("", "access_token", "/"+u.Name, time.Now().Add(-1*time.Hour))
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
