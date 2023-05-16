package entity

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID
	Username string
	Password string
}

type Admin struct {
	User
}

type Customer struct {
	User
	Balance
}

type Balance struct {
	total map[string]float64
}
