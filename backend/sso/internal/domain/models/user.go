package models

type User struct {
	ID       uint64
	Login    string
	Email    string
	PassHash []byte
}
