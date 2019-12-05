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

func (cpr *CompanyRepoImpl) GetAll() ([]model.Company, error) {
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

func (cpr *CompanyRepoImpl) GetCompaniesWithUsers() ([]model.CompanyWithUsers, error) {

	type CompanyWithUsersSet map[int]*model.CompanyWithUsers
	companies := make(CompanyWithUsersSet, 0)

	rows, err := cpr.dbConn.Query("select iu.iam_id,iu.company_id,c.\"name\",c.duns,c.spin,c.city,c.country,c.\"type\" from iam_user iu left join company c on c.id=iu.company_id")
	if err != nil {
		return make([]model.CompanyWithUsers, 0), err
	}

	for rows.Next() {
		company := model.Company{}
		iamID := ""
		if err := rows.Scan(&iamID, &company.ID, &company.Name, &company.Duns, &company.Spin, &company.City, &company.Country, &company.Type); err != nil {
			return make([]model.CompanyWithUsers, 0), err
		}

		if companies[company.ID] == nil {
			companies[company.ID] = &model.CompanyWithUsers{company, make([]string, 0)}
		}
		companies[company.ID].IamIDs = append(companies[company.ID].IamIDs, iamID)
	}

	result := make([]model.CompanyWithUsers, 0)
	for _, company := range companies {
		result = append(result, *company)
	}

	return result, err
}
