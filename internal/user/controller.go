package user

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/civil"
	"github.com/ihopethisisfine/helloworld/internal/domain"
	"github.com/ihopethisisfine/helloworld/internal/pkg/storage"
)

type Controller struct {
	Storage storage.UserStorer
}

func (c Controller) Hello(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		c.find(w, r)
	case http.MethodPut:
		c.put(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (c Controller) put(w http.ResponseWriter, r *http.Request) {
	var req User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid JSON")
		return
	}

	user := storage.User{
		Username:    strings.TrimPrefix(r.URL.Path, "/hello/"),
		DateOfBirth: req.DateOfBirth,
	}

	err := user.Validate()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	err = c.Storage.Put(r.Context(), user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	_, _ = fmt.Fprint(w, "")
}

func (c Controller) find(w http.ResponseWriter, r *http.Request) {
	res, err := c.Storage.Find(r.Context(), strings.TrimPrefix(r.URL.Path, "/hello/"))
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(Response{Message: err.Error()})
		return
	}
	data := Response{Message: ""}
	daysUntilBirthday := birthdayCountdown(res.DateOfBirth)

	if daysUntilBirthday == 0 {
		data.Message = fmt.Sprintf("Hello, %s! Happy birthday!", res.Username)
	} else {
		data.Message = fmt.Sprintf("Hello, %s! Your birthday is in %d day(s)", res.Username, daysUntilBirthday)
	}

	_ = json.NewEncoder(w).Encode(data)
}

// Returns days until next birthdate
func birthdayCountdown(birthdate string) int {
	parsedBirthdate, _ := time.Parse("2006-01-02", birthdate)
	birthday := civil.DateOf(parsedBirthdate)
	today := civil.DateOf(time.Now())
	birthday.Year = today.Year
	if birthday.Day == 29 && birthday.Month == 2 {
		if !isLeap(birthday) {
			birthday.Month = 3
			birthday.Day = 1
		}
	}

	days := today.DaysSince(birthday)
	if days > 0 {
		birthday.Year = birthday.Year + 1
		days = today.DaysSince(birthday)
	}
	days = days * -1
	return days
}

// Returns true if year date is a leap year.
func isLeap(date civil.Date) bool {
	year := date.Year
	if year%400 == 0 {
		return true
	} else if year%100 == 0 {
		return false
	} else if year%4 == 0 {
		return true
	}
	return false
}
