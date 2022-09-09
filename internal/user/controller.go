package user

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/ihopethisisfine/helloworld/internal/domain"
	"github.com/ihopethisisfine/helloworld/internal/pkg/storage"

	"github.com/google/uuid"
)

type Controller struct {
	Storage storage.UserStorer
}

func (c Controller) Hello(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		c.find(w, r)
	case http.MethodPut:
		c.create(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (c Controller) create(w http.ResponseWriter, r *http.Request) {
	var req User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id := uuid.New().String()

	err := c.Storage.Insert(r.Context(), storage.User{
		Username:    strings.TrimPrefix(r.URL.Path, "/hello/"),
		DateOfBirth: req.DateOfBirth,
	})
	if err != nil {
		switch err {
		case domain.ErrConflict:
			w.WriteHeader(http.StatusConflict)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(id))
}

func (c Controller) find(w http.ResponseWriter, r *http.Request) {
	res, err := c.Storage.Find(r.Context(), strings.TrimPrefix(r.URL.Path, "/hello/"))
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	user := User{
		DateOfBirth: res.DateOfBirth,
	}

	data, err := json.Marshal(user)
	if err != nil {
		log.Println(err)

		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = w.Write(data)
}
