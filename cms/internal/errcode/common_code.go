package errcode

import "data_backend/pkg/errcode"

var (
	// User
	IncorrectPassword = errcode.ErrorParam.WithMsg("IncorrectPassword")
	UserIsLock        = errcode.Forbidden.WithMsg("UserIsLock")
	LoginFail         = errcode.ServerError.WithMsg("LoginFail")

	// Exist
	UserExist    = errcode.NotAcceptable.WithMsg("UserExist")
	UserNotExist = errcode.NotAcceptable.WithMsg("UserNotExist")
	RoleExist    = errcode.NotAcceptable.WithMsg("RoleExist")
	RoleNotExist = errcode.NotAcceptable.WithMsg("RoleNotExist")

	// Menu
	InitMenuFail = errcode.ServerError.WithMsg("InitMenuFail")

	// SQl
	SQLExecFail = errcode.ServerError.WithMsg("SQLExecFail")
)
