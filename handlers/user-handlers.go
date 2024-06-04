package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/saalikmubeen/goravel"
	"github.com/saalikmubeen/goravel-demo-app/models"
)

func (h *Handlers) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.Models.Users.GetAll()
	if err != nil {
		h.App.ErrorLog.Println(err)
		return
	}

	h.App.WriteJSON(w, 200, goravel.Response{
		"users": users,
	})
}

func (h *Handlers) GetUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	u, err := h.Models.Users.Get(id)
	if err != nil {
		h.App.WriteJSON(w, 400, goravel.Response{
			"error": err.Error(),
		})
		return
	}

	h.App.WriteJSON(w, 200, goravel.Response{
		"user": u,
	})
}

func (h *Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	u := models.User{
		FirstName: "Colt",
		LastName:  "Steele",
		Email:     "colt@here.com",
		Active:    1,
		Password:  "password",
	}

	validation := h.App.Validator(nil)

	u.Validate(validation)

	if !validation.IsValid() {
		h.App.WriteJSON(w, 400, goravel.Response{
			"errors": validation.Errors,
		})
		return

	}

	id, err := h.Models.Users.Insert(u)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return
	}

	u.ID = id

	h.App.WriteJSON(w, 200, goravel.Response{
		"user": u,
	})
}

func (h *Handlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	u, err := h.Models.Users.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return
	}

	u.LastName = h.App.RandomString(10)

	validator := h.App.Validator(nil)
	u.LastName = ""

	u.Validate(validator)

	if !validator.IsValid() {
		fmt.Fprint(w, "failed validation")
		return
	}
	err = u.Update(*u)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return
	}

	fmt.Fprintf(w, "updated last name to %s", u.LastName)
}

func (h *Handlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	u, err := h.Models.Users.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return
	}

	err = u.Delete(u.ID)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return
	}

	fmt.Fprintf(w, "deleted user %d", id)
}
