package rest

import (
	"html/template"
	"net/http"
	"os"
	"path"
	"strings"
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

	w.WriteHeader(http.StatusOK)
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

	w.WriteHeader(http.StatusOK)
}

func (h Handler) RegisterTemplate(w http.ResponseWriter, r *http.Request) {
	fp := path.Join(h.PathToTemplates()+"/templates/mainmenu", "registerform.html")

	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		h.log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "index", nil); err != nil {
		h.log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h Handler) CreateTasksTemplate(w http.ResponseWriter, r *http.Request) {
	fp := path.Join(h.PathToTemplates()+"/templates/mainmenu", "createtask.html")

	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		h.log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "index", nil); err != nil {
		h.log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
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
