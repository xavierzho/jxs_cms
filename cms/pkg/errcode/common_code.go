package errcode

var (
	// 400
	BadRequest    = newError(400, 40000000, "BadRequest")    // 400
	InvalidParams = newError(400, 40000001, "InvalidParams") // 输入的参数不符合要求
	ErrorParam    = newError(400, 40000002, "ErrorParam")    // 输入的信息有误

	// 401
	Unauthorized             = newError(401, 40100000, "Unauthorized")
	UnauthorizedTokenError   = newError(401, 40100001, "UnauthorizedTokenError")   // Token 信息错误
	UnauthorizedTokenTimeout = newError(401, 40100002, "UnauthorizedTokenTimeout") // Token 超时
	UnauthorizedTokenOverdue = newError(401, 40100003, "UnauthorizedTokenOverdue") // 因超时之外的原因过期--修改了密码, 在其他地方重新登录

	Forbidden        = newError(403, 40300000, "Forbidden")
	TooManyRequests  = newError(403, 40300001, "TooManyRequests")
	PermissionDenied = newError(403, 40300002, "PermissionDenied")

	NotFound = newError(404, 40400000, "NotFound")

	NotAcceptable = newError(406, 40600000, "NotAcceptable") // 无法执行用户请求--无法重复创建, 无有效目标

	// 500
	ServerError    = newError(500, 50000000, "ServerError")
	ExecuteFail    = newError(500, 50000001, "ExecuteFail")
	CreateFail     = newError(500, 50000002, "CreateFail")
	DeleteFail     = newError(500, 50000003, "DeleteFail")
	QueryFail      = newError(500, 50000004, "QueryFail")
	UpdateFail     = newError(500, 50000005, "UpdateFail")
	TransformFail  = newError(500, 50000011, "TransformFail")
	ExportFail     = newError(500, 50000012, "ExportFail")
	UploadFileFail = newError(500, 50000013, "UploadFileFail")
)
