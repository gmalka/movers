package rest

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/gmalka/movers/model"
)

var pathToTemplates string

func (h Handler) MainMenu(w http.ResponseWriter, r *http.Request) {
	fp := path.Join(h.PathToTemplates()+"/templates/mainmenu", "mainmenu.html")

	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		h.log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "main", nil); err != nil {
		h.log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h Handler) LoginTemplate(w http.ResponseWriter, r *http.Request) {
	fp := path.Join(h.PathToTemplates()+"/templates/mainmenu", "loginform.html")

	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		h.log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "login", nil); err != nil {
		h.log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	u := model.User{}

	err := r.ParseForm()
	if err != nil {
		h.log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	u.Name = r.Form.Get("login")
	u.Password = r.Form.Get("password")

	fmt.Println(u)
}

func (h Handler) RegisterTemplate(w http.ResponseWriter, r *http.Request) {
	fp := path.Join(h.PathToTemplates()+"/templates/mainmenu", "registerform.html")

	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		h.log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "register", nil); err != nil {
		h.log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h Handler) Regsiter(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	user := model.User{}
	user.Name = r.Form.Get("login")
	user.Password = r.Form.Get("password")
	user.Role = r.Form.Get("options")
	// err = h.auth.Register(r.Context(), user)
	// if err != nil {
	// 	h.log.Error(err.Error())
	// 	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	// 	return
	// }

	switch user.Role {
	case "Customer":
		customer := model.CustomerInfo{}
		money, err := strconv.Atoi(r.Form.Get("money"))
		if err != nil {
			//h.auth.DeleteUser(r.Context(), user.Name)
			h.log.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		customer.Money = money
		customer.Name = user.Name

		// err = h.users.NewCustomer(r.Context(), customer)
		// if err != nil {
		// 	h.auth.DeleteUser(r.Context(), user.Name)
		// 	h.log.Error(err.Error())
		// 	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		// 	return
		// }
		fmt.Println(customer)
	case "Worker":
		worker := model.WorkerInfo{}
		fatigue, err := strconv.Atoi(r.Form.Get("fatigue"))
		if err != nil {
			//h.auth.DeleteUser(r.Context(), user.Name)
			h.log.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		worker.Fatigue = fatigue

		salary, err := strconv.Atoi(r.Form.Get("price"))
		if err != nil {
			//h.auth.DeleteUser(r.Context(), user.Name)
			h.log.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		worker.Salary = salary

		weight, err := strconv.Atoi(r.Form.Get("weight"))
		if err != nil {
			//h.auth.DeleteUser(r.Context(), user.Name)
			h.log.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		worker.CarryWeight = weight

		drunk := r.Form.Get("drunk")
		if drunk == "" {
			worker.Drunk = 1
		} else {
			worker.Drunk = 2
		}
	default:
		//h.auth.DeleteUser(r.Context(), user.Name)
		h.log.Error(fmt.Sprintf("Incorrect role: %s", user.Role))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}

func (h Handler) CreateTasksTemplate(w http.ResponseWriter, r *http.Request) {
	fp := path.Join(h.PathToTemplates()+"/templates/mainmenu", "createtask.html")

	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		h.log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "tasks", nil); err != nil {
		h.log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
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
