package postgres

import (
	"fmt"

	"github.com/thomasobenaus/goms/model"
)

func (cpr *PGRepo) Add(user model.User) error {

	queryStr := fmt.Sprintf("insert into iam_user (iam_id, email, name, company_id) "+
		"values('%s','%s','%s',%d)", user.IamID, user.Email, user.Name, user.CompanyID)

	_, err := cpr.dbConn.Query(queryStr)
	if err != nil {
		return err
	}

	return nil
}
