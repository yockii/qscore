package constant

const (
	ErrorCodeUnknown = -10000 - iota
	ErrorCodeBodyParse
	ErrorCodeLackOfField
	ErrorCodeNotFound
	ErrorCodeService
	ErrorCodeDuplicate
	ErrorCodeInvalid
	ErrorCodeReject
)
