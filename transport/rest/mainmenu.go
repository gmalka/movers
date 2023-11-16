package rest

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gmalka/movers/model"
)

var pathToTemplates string

func (h Handler) MainMenu(w http.ResponseWriter, r *http.Request) {
	fp := path.Join(h.PathToTemplates()+"/templates/mainmenu", "mainmenu.html")

	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		h.HandlerError(w, err, http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "main", nil); err != nil {
		h.HandlerError(w, err, http.StatusInternalServerError)
		return
	}
}

func (h Handler) LoginTemplate(w http.ResponseWriter, r *http.Request) {
	fp := path.Join(h.PathToTemplates()+"/templates/mainmenu", "loginform.html")

	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		h.HandlerError(w, err, http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "login", nil); err != nil {
		h.HandlerError(w, err, http.StatusInternalServerError)
		return
	}
}

func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	u := model.User{}

	err := r.ParseForm()
	if err != nil {
		h.HandlerError(w, err, http.StatusInternalServerError)
		return
	}

	u.Name = r.Form.Get("login")
	u.Password = r.Form.Get("password")

	tokens, err := h.auth.Login(r.Context(), u.Name, u.Password)
	if err != nil {
		h.HandlerError(w, err, http.StatusInternalServerError)
		return
	}

	cookie := CreateBearerTokenCookie(tokens.AccessToken, "access_token", "/" + u.Name, time.Now().Add(h.auth.GetAccessTTL()*time.Minute))
	http.SetCookie(w, cookie)

	b, err := json.Marshal(tokens)
	if err != nil {
		h.HandlerError(w, err, http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

func CreateBearerTokenCookie(token, tokenname, path string, expiration time.Time) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = tokenname
	cookie.Value = "Bearer " + token
	cookie.Expires = expiration
	cookie.Path = path
	cookie.HttpOnly = true

	return cookie
}

func (h Handler) RegisterTemplate(w http.ResponseWriter, r *http.Request) {
	fp := path.Join(h.PathToTemplates()+"/templates/mainmenu", "registerform.html")

	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		h.HandlerError(w, err, http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "register", nil); err != nil {
		h.HandlerError(w, err, http.StatusInternalServerError)
		return
	}
}

func (h Handler) Regsiter(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.HandlerError(w, err, http.StatusInternalServerError)
		return
	}

	user := model.User{}
	user.Name = r.Form.Get("login")
	user.Password = r.Form.Get("password")
	user.Role = r.Form.Get("options")
	h.log.Info(fmt.Sprintf("Createing user %s", user.Name))
	err = h.auth.Register(r.Context(), user)
	if err != nil {
		h.HandlerError(w, err, http.StatusInternalServerError)
		return
	}

	switch user.Role {
	case "Customer":
		customer := model.CustomerInfo{}
		money, err := strconv.Atoi(r.Form.Get("money"))
		if err != nil {
			h.auth.DeleteUser(r.Context(), user.Name)
			h.HandlerError(w, err, http.StatusInternalServerError)
			return
		}
		customer.Money = money
		customer.Name = user.Name

		h.log.Info(fmt.Sprintf("Createing Customer %#v", customer))
		err = h.users.NewCustomer(r.Context(), customer)
		if err != nil {
			h.auth.DeleteUser(r.Context(), user.Name)
			h.HandlerError(w, err, http.StatusInternalServerError)
			return
		}
	case "Worker":
		worker := model.WorkerInfo{Name: user.Name}
		fatigue, err := strconv.Atoi(r.Form.Get("fatigue"))
		if err != nil {
			h.auth.DeleteUser(r.Context(), user.Name)
			h.HandlerError(w, err, http.StatusInternalServerError)
			return
		}
		worker.Fatigue = fatigue

		salary, err := strconv.Atoi(r.Form.Get("price"))
		if err != nil {
			h.auth.DeleteUser(r.Context(), user.Name)
			h.HandlerError(w, err, http.StatusInternalServerError)
			return
		}
		worker.Salary = salary

		weight, err := strconv.Atoi(r.Form.Get("weight"))
		if err != nil {
			h.auth.DeleteUser(r.Context(), user.Name)
			h.HandlerError(w, err, http.StatusInternalServerError)
			return
		}
		worker.CarryWeight = weight

		drunk := r.Form.Get("drunk")
		if drunk == "" {
			worker.Drunk = 1
		} else {
			worker.Drunk = 2
		}

		h.log.Info(fmt.Sprintf("Createing Worker %#v", worker))
		err = h.users.NewWorker(r.Context(), worker)
		if err != nil {
			h.auth.DeleteUser(r.Context(), user.Name)
			h.HandlerError(w, err, http.StatusInternalServerError)
			return
		}
	default:
		h.auth.DeleteUser(r.Context(), user.Name)
		h.HandlerError(w, fmt.Errorf("incorrect role: %s", user.Role), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h Handler) CreateTasksTemplate(w http.ResponseWriter, r *http.Request) {
	fp := path.Join(h.PathToTemplates()+"/templates/mainmenu", "createtask.html")

	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		h.HandlerError(w, err, http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "tasks", nil); err != nil {
		h.HandlerError(w, err, http.StatusInternalServerError)
		return
	}
}

func (h Handler) CreateTasks(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.HandlerError(w, err, http.StatusInternalServerError)
		return
	}

	count, err := strconv.Atoi(r.Form.Get("tasks"))
	if err != nil {
		h.HandlerError(w, err, http.StatusInternalServerError)
		return
	}

	err = h.tasks.GenerateTasks(r.Context(), count)
	if err != nil {
		h.HandlerError(w, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h Handler) PathToTemplates() string {
	var err error

	if pathToTemplates == "" {
		pathToTemplates, err = os.Getwd()
		if err != nil {
			h.log.Error("NO PWD")
			return pathToTemplates
		}
		strs := strings.Split(pathToTemplates, "/")
		str := ""
		for i := range strs {
			if strs[i] == "my_site" {
				str += "/" + strs[i]
				break
			}
			str += "/" + strs[i]
		}

		pathToTemplates = str
	}

	return pathToTemplates
}
