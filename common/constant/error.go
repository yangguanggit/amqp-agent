package constant

const (
	ErrorOk          = 200 //正确
	ErrorSystemError = 500 //系统异常
	ErrorParams      = 400 //参数错误
)

var errorMap = map[int]string{
	ErrorOk:          "成功",
	ErrorSystemError: "系统异常，请稍后重试",
	ErrorParams:      "参数错误",
}

/**
 * 获取错误信息
 */
func ErrorMessage(code int) string {
	if message, ok := errorMap[code]; ok {
		return message
	}
	return "系统异常，请稍后重试"
}
