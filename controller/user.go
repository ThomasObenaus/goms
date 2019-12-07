package controller

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/thomasobenaus/goms/model"
)

type UserController struct {
	repo model.UserRepo
}

func NewUserController(userRepo model.UserRepo) UserController {
	return UserController{userRepo}
}

func (uco *UserController) AddUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	user := model.User{}

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := uco.repo.Add(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(w)
	if err := enc.Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
