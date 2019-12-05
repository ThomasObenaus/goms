package model

type Company struct {
	ID      int    `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Duns    string `json:"duns,omitempty"`
	Spin    string `json:"spin,omitempty"`
	City    string `json:"city,omitempty"`
	Country string `json:"country,omitempty"`
	Type    string `json:"type,omitempty"`
}

type CompanyWithUsers struct {
	Company Company  `json:"company,omitempty"`
	IamIDs  []string `json:"users,omitempty"`
}

type CompanyRepo interface {
	Get(id int) (Company, error)
	GetAll() ([]Company, error)
	GetCompaniesWithUsers() ([]CompanyWithUsers, error)
	//FindBy(Type string) ([]Company, error)
	//Delete(company Company)
	//Create(company Company)
	//Update(company Company)
}
