package constant

const (
	ResponseCodeErrorInternal = -iota - 10000
	ResponseCodeErrorParamParse
	ResponseCodeErrorParamRequired
	ResponseCodeErrorDataExists
	ResponseCodeErrorUnsupported
	ResponseCodeErrorResourceNotFound
	ResponseCodeErrorDataNotMatch
	ResponseCodeErrorDataInvalid
	ResponseCodeErrorDataLimitation
)

const (
	ResponseMsgErrorInternal         = "service error"
	ResponseMsgErrorParamParse       = "failed to parse params"
	ResponseMsgErrorParamRequired    = "param is required"
	ResponseMsgErrorDataExists       = "data exists"
	ResponseMsgErrorUnsupported      = "unsupported"
	ResponseMsgErrorResourceNotFound = "resource not found"
	ResponseMsgErrorDataNotMatch     = "data not match"
	ResponseMsgErrorDataInvalid      = "data invalid"
	ResponseMsgErrorDataLimitation   = "data limitation reached"
)
