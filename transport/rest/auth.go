package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
)

func (h Handler) checkAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenRaw string
		cookie, err := r.Cookie("access_token")
		if err == nil {
			tokenRaw = cookie.Value
		} else {
			tokenRaw = r.Header.Get("Authorization")
		}

		tokenParts := strings.Split(tokenRaw, " ")
		if len(tokenParts) < 2 && tokenParts[0] != "Bearer" {
			h.log.Error(fmt.Sprintf("wrong authorization: %v\n", tokenParts))
			http.Error(w, "message: wrong authorization token", http.StatusBadRequest)
			return
		}

		u, err := h.auth.CheckAccessToken(tokenParts[1])
		if err != nil {
			h.log.Error(err.Error())
			http.Error(w, "message: wrong authorization token", http.StatusBadRequest)
			return
		}

		username := chi.URLParam(r, "username")
		if username != u.Name {
			h.log.Error(fmt.Sprintf("username in token and path are different: %s-%s", username, u.Name))
			http.Error(w, "message: invalid resource", http.StatusNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), UserRequest{}, u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h Handler) refreshAccessToken(w http.ResponseWriter, r *http.Request) {
	var err error

	refresh := struct{
		RefreshToken string `json:"refresh_token"`
	}{}

	access := struct{
		AccessToken string `json:"access_token"`
	}{}

	access.AccessToken, err = h.auth.UpdateAccessToken(refresh.RefreshToken)
	if err != nil {
		h.HandlerError(w, err, http.StatusBadRequest)
	}

	b, _ := json.Marshal(access)
	w.Write(b)
}

func (h Handler) refreshRefreshToken(w http.ResponseWriter, r *http.Request) {
	var err error

	refresh := struct{
		RefreshToken string `json:"refresh_token"`
	}{}

	refresh.RefreshToken, err = h.auth.UpdateRefreshToken(refresh.RefreshToken)
	if err != nil {
		h.HandlerError(w, err, http.StatusBadRequest)
	}

	b, _ := json.Marshal(refresh)
	w.Write(b)
}
