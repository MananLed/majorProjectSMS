package utils

type contextKey string

const (
	UserIDKey     contextKey = "user_id"
	UserEmailKey  contextKey = "user_email"
	UserRoleKey   contextKey = "user_role"
	UserFlatKey   contextKey = "user_flat"
)