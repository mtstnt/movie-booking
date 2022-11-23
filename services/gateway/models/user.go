package models

type User struct {
	ID       uint32
	Email    string
	Name     string
	Password string `json:"-"`
}
