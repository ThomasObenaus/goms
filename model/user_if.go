package model

type User struct {
	IamID   string
	Name    string
	Surname string
}

type UserRepo interface {
	Add(user User) error
}
