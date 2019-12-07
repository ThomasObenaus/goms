package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/schema"
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
		http.Error(w, err.Error(), http.StatusBadRequest)
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

type Paginate struct {
	PageSize int
	Page     int
}

func (coc *CompanyController) GetCompaniesWithUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	code := http.StatusOK

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	paginate := Paginate{}
	decoder := schema.NewDecoder()
	if err := decoder.Decode(&paginate, r.Form); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	companies, totalPages, totalElements, err := coc.repo.GetCompaniesWithUsers(paginate.Page, paginate.PageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("TotalPages", strconv.Itoa(totalPages))
	w.Header().Add("TotalElements", strconv.Itoa(totalElements))

	w.WriteHeader(code)

	enc := json.NewEncoder(w)
	if err := enc.Encode(companies); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
