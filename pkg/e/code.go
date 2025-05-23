package e

const (
	SUCCESS               = 200
	UpdatePasswordSuccess = 201
	NotExistInentifier    = 202
	ERROR                 = 500
	InvalidParams         = 400

	ErrUserNotFound     = 10001
	ErrPasswordMismatch = 10002
	ErrUsernameExists   = 10003

	//数据库错误
	ErrorDatabase          = 40000
	ErrDatabaseQueryFailed = 40001
)
