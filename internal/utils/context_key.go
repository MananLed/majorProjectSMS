package utils

type contextKey string

const (
	UserIDKey     contextKey = "user_id"
	UserRoleKey   contextKey = "user_role"
	UserPassKey   contextKey = "user_password"
)