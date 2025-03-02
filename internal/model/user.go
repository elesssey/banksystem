package model

type Role string

const (
	RoleClient     Role = "client"
	RoleOperator   Role = "operator"
	RoleManager    Role = "manager"
	RoleSpecialist Role = "specialist"
	RoleAdmin      Role = "admin"
)

type User struct {
	ID             int
	Username       string
	Name           string
	Surname        string
	PassportSeries string
	PassportNumber string
	Phone          string
	Email          string
	Role           Role
}
