package user

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/civil"
	"github.com/ihopethisisfine/helloworld/internal/domain"
	"github.com/ihopethisisfine/helloworld/internal/pkg/storage"
)

var _ storage.UserStorer = UserStorageMock{}

var useCurrentDateForBirthdate = false

type UserStorageMock struct {
}

var usr = Controller{
	Storage: UserStorageMock{},
}

func (u UserStorageMock) Put(ctx context.Context, user storage.User) error {
	return nil
}

func (u UserStorageMock) Find(ctx context.Context, username string) (storage.User, error) {
	birthdate := time.Now().Format("2006-01-02")
	if !useCurrentDateForBirthdate {
		tomorrow := time.Now().AddDate(0, 0, 1)
		birthdate = tomorrow.Format("2006-01-02")
	}
	return storage.User{Username: "Mike", DateOfBirth: birthdate}, nil
}

func TestHello(t *testing.T) {
	t.Run("can register a valid user", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodPut, "/hello/Mike", strings.NewReader("{\"dateOfBirth\": \"2000-02-02\"}"))
		res := httptest.NewRecorder()

		usr.Hello(res, req)

		if res.Code != http.StatusNoContent {
			t.Errorf("expected status of %d but got %d", http.StatusNoContent, res.Code)
		}
	})

	t.Run("returns 400 bad request if username is invalid", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodPut, "/hello/Mike123", strings.NewReader("{\"dateOfBirth\": \"2000-02-02\"}"))
		res := httptest.NewRecorder()

		usr.Hello(res, req)

		if res.Code != http.StatusBadRequest {
			t.Errorf("expected status of %d but got %d", http.StatusBadRequest, res.Code)
		}

		if res.Body.String() != domain.ErrInvalidUsername.Error() {
			t.Errorf("expected body of %s but got %s", domain.ErrInvalidUsername.Error(), res.Body.String())
		}
	})

	t.Run("returns 400 bad request if date is invalid", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodPut, "/hello/Mike", strings.NewReader("{\"dateOfBirth\": \"2000-02-31\"}"))
		res := httptest.NewRecorder()

		usr.Hello(res, req)

		if res.Code != http.StatusBadRequest {
			t.Errorf("expected status of %d but got %d", http.StatusBadRequest, res.Code)
		}

		if res.Body.String() != domain.ErrInvalidDate.Error() {
			t.Errorf("expected body of %s but got %s", domain.ErrInvalidDate.Error(), res.Body.String())
		}
	})

	t.Run("returns 400 bad request if date is not in format YYYY-MM-DD", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodPut, "/hello/Mike", strings.NewReader("{\"dateOfBirth\": \"20-02-2002\"}"))
		res := httptest.NewRecorder()

		usr.Hello(res, req)

		if res.Code != http.StatusBadRequest {
			t.Errorf("expected status of %d but got %d", http.StatusBadRequest, res.Code)
		}

		if res.Body.String() != domain.ErrInvalidDate.Error() {
			t.Errorf("expected body of %s but got %s", domain.ErrInvalidDate.Error(), res.Body.String())
		}
	})

	t.Run("returns 400 bad request if birthdate is after the current date", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodPut, "/hello/Mike", strings.NewReader("{\"dateOfBirth\": \"9999-02-01\"}"))
		res := httptest.NewRecorder()

		usr.Hello(res, req)

		if res.Code != http.StatusBadRequest {
			t.Errorf("expected status of %d but got %d", http.StatusBadRequest, res.Code)
		}

		if res.Body.String() != domain.ErrInvalidBirthDate.Error() {
			t.Errorf("expected body of %s but got %s", domain.ErrInvalidBirthDate.Error(), res.Body.String())
		}
	})

	t.Run("returns 400 bad request if is not valid user JSON", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodPut, "/hello/Mike", strings.NewReader("not json"))
		res := httptest.NewRecorder()

		usr.Hello(res, req)

		if res.Code != http.StatusBadRequest {
			t.Errorf("expected status of %d but got %d", http.StatusBadRequest, res.Code)
		}

		if res.Body.String() != "Invalid JSON" {
			t.Errorf("expected body of %s but got %s", "Invalid JSON", res.Body.String())
		}
	})

	t.Run("return how many days until user's birthday", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodGet, "/hello/Mike", nil)
		res := httptest.NewRecorder()

		useCurrentDateForBirthdate = false
		usr.Hello(res, req)

		if res.Code != http.StatusOK {
			t.Errorf("expected status of %d but got %d", http.StatusOK, res.Code)
		}

		expectedBody := Response{Message: "Hello, Mike! Your birthday is in 1 day(s)"}
		response := decodeResponse(res.Body)

		if response.Message != expectedBody.Message {
			t.Errorf("expected body of %s but got %s", expectedBody, response)
		}
	})

	t.Run("wish user happy birthday if it is today", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodGet, "/hello/Mike", strings.NewReader("not json"))
		res := httptest.NewRecorder()

		useCurrentDateForBirthdate = true

		usr.Hello(res, req)

		if res.Code != http.StatusOK {
			t.Errorf("expected status of %d but got %d", http.StatusNoContent, res.Code)
		}

		expectedBody := Response{Message: "Hello, Mike! Happy birthday!"}
		response := decodeResponse(res.Body)

		if response.Message != expectedBody.Message {
			t.Errorf("expected body of %s but got %s", expectedBody, response)
		}
	})
}

func decodeResponse(body *bytes.Buffer) Response {
	var r Response
	if err := json.NewDecoder(body).Decode(&r); err != nil {
		log.Print(err)
	}
	return r
}

func Test_birthdayCountdown(t *testing.T) {
	type args struct {
		birthdate string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"Should return 10 if birthday is today+10",
			args{birthdate: time.Now().AddDate(0, 0, 10).Format("2006-01-02")},
			10,
		},
		{
			"Should return a whole year if birthday was yesterday",
			args{birthdate: time.Now().AddDate(0, 0, -1).Format("2006-01-02")},
			civil.DateOf(time.Now().AddDate(1, 0, 0)).DaysSince(civil.DateOf(time.Now().AddDate(0, 0, 1))),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := birthdayCountdown(tt.args.birthdate); got != tt.want {
				t.Errorf("birthdayCountdown() = %v, want %v", got, tt.want)
			}
		})
	}
}
