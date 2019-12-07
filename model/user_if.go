package model

type User struct {
	IamID     string
	Email     string
	Name      string
	CompanyID int
}

type UserRepo interface {
	Add(user User) error
}
