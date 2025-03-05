package model

type Bank struct {
	ID          int
	Name        string
	Descrition  string
	BIC         string
	Address     string
	Rating      int
	Enterprises []*Enterprise
}
