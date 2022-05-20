package model

type AuthType string

const (
	Password = AuthType("password")
	Google   = AuthType("google")
	Github   = AuthType("github")
)
