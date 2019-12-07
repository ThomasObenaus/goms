package postgres

import (
	"math"

	"github.com/thomasobenaus/goms/model"
)

type PGRepo struct {
	dbConn *DBConnection
}

func NewPGRepo(dbConn *DBConnection) *PGRepo {

	return &PGRepo{dbConn: dbConn}
}

func (cpr *PGRepo) Get(id int) (model.Company, error) {

	result := model.Company{}

	rows, err := cpr.dbConn.Query("select * from company c where c.id=$1;", id)
	if err != nil {
		return result, err
	}

	if rows.Next() {
		if err := rows.Scan(&result.ID, &result.Name, &result.Duns, &result.Spin, &result.City, &result.Country, &result.Type); err != nil {
			return model.Company{}, err
		}
	}

	return result, nil
}

func (cpr *PGRepo) GetAll() ([]model.Company, error) {
	result := make([]model.Company, 0)

	rows, err := cpr.dbConn.Query("select * from company")
	if err != nil {
		return result, err
	}

	for rows.Next() {
		company := model.Company{}
		if err := rows.Scan(&company.ID, &company.Name, &company.Duns, &company.Spin, &company.City, &company.Country, &company.Type); err != nil {
			return make([]model.Company, 0), err
		}

		result = append(result, company)
	}

	return result, nil
}

func (cpr *PGRepo) GetCompaniesWithUsers(page, pageSize int) (companies []model.CompanyWithUsers, totalPages int, totalElements int, err error) {

	type CompanyWithUsersSet map[int]*model.CompanyWithUsers
	companySet := make(CompanyWithUsersSet, 0)

	rows, err := cpr.dbConn.Query("select iu.iam_id,iu.company_id,c.\"name\",c.duns,c.spin,c.city,c.country,c.\"type\" from iam_user iu left join company c on c.id=iu.company_id")
	if err != nil {
		return make([]model.CompanyWithUsers, 0), 0, 0, err
	}

	for rows.Next() {
		company := model.Company{}
		iamID := ""
		if err := rows.Scan(&iamID, &company.ID, &company.Name, &company.Duns, &company.Spin, &company.City, &company.Country, &company.Type); err != nil {
			return make([]model.CompanyWithUsers, 0), 0, 0, err
		}

		if companySet[company.ID] == nil {
			companySet[company.ID] = &model.CompanyWithUsers{company, make([]string, 0)}
		}
		companySet[company.ID].IamIDs = append(companySet[company.ID].IamIDs, iamID)
	}

	result := make([]model.CompanyWithUsers, 0)
	for _, company := range companySet {
		result = append(result, *company)
	}

	numTotalElements := len(result)
	numTotalPages := math.Ceil(float64(numTotalElements) / float64(pageSize))

	start := page * pageSize
	end := (page + 1) * pageSize
	if start > numTotalElements {
		start = numTotalElements
	}

	if end > numTotalElements {
		end = numTotalElements
	}

	result = result[start:end]
	return result, int(numTotalPages), numTotalElements, err
}
