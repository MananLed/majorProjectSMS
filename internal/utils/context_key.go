package utils

type contextKey string

const (
	UserIDKey     contextKey = "user_id"
	UserEmailKey  contextKey = "user_email"
	UserRoleKey   contextKey = "user_role"
	UserPassKey   contextKey = "user_password"
	UserFlatKey   contextKey = "user_flat"
)