package server

const (
	ResponseCodeUnknownError = -10000 - iota
	ResponseCodeParamParseError
	ResponseCodeParamNotEnough
	ResponseCodePasswordStrengthInvalid
	ResponseCodeDuplicated
	ResponseCodeDatabase
	ResponseCodeDataNotExists
	ResponseCodeDataNotMatch
	ResponseCodeModuleNotExists
	ResponseCodeGeneration
)

var (
	ResponseMsgUnknownError            = "系统错误"
	ResponseMsgParamParseError         = "参数解析失败"
	ResponseMsgParamNotEnough          = "参数不足"
	ResponseMsgPasswordStrengthInvalid = "密码强度不够"
	ResponseMsgDuplicated              = "数据重复"
	ResponseMsgDatabase                = "执行数据库语句失败"
	ResponseMsgDataNotExists           = "数据不存在"
	ResponseMsgDataNotMatch            = "数据不匹配"
	ResponseMsgModuleNotExists         = "模块不存在"
	ResponseMsgGeneration              = "生成信息失败"
)
