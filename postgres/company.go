package postgres

import "github.com/thomasobenaus/goms/model"

type CompanyRepoImpl struct {
	dbConn *DBConnection
}

func NewPGCompanyRepo(dbConn *DBConnection) model.CompanyRepo {

	return &CompanyRepoImpl{dbConn: dbConn}
}

func (cpr *CompanyRepoImpl) Get(id int) (model.Company, error) {

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
