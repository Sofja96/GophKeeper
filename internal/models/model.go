package models

type ContextKey string

const ContextKeyUser ContextKey = "username"

type User struct {
	Username string `db:"username"`
	Password string `db:"password"`
}
