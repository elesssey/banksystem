package model

type Role string

const (
	RoleClient   Role = "client"
	RoleOperator Role = "operator"
	RoleManager  Role = "manager"
	RoleAdmin    Role = "admin"
)

type User struct {
	ID             int
	Password       string
	Name           string
	MiddleName     string
	Surname        string
	PassportSeries string
	PassportNumber string
	Phone          string
	Email          string
	Role           Role
}
