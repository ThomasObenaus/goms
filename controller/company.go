package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/thomasobenaus/goms/model"
)

type CompanyController struct {
	repo model.CompanyRepo
}

func New(companyRepo model.CompanyRepo) CompanyController {
	return CompanyController{companyRepo}
}

func (coc *CompanyController) GetCompany(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	code := http.StatusOK

	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	company, err := coc.repo.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)

	enc := json.NewEncoder(w)
	if err := enc.Encode(company); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (coc *CompanyController) GetCompanies(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	code := http.StatusOK

	companies, err := coc.repo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)

	enc := json.NewEncoder(w)
	if err := enc.Encode(companies); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (coc *CompanyController) GetCompaniesWithUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	code := http.StatusOK

	companies, err := coc.repo.GetCompaniesWithUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)

	enc := json.NewEncoder(w)
	if err := enc.Encode(companies); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
