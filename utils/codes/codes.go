package codes

// 统一业务状态码定义。
//
// 说明：
// 1) code=0 代表成功；
// 2) 4xxxx 代表参数或调用方错误；
// 3) 5xxxx 代表网关/依赖服务错误；
// 4) 各服务新增码值时优先在此集中维护，避免散落硬编码。
const (
	// OK 请求成功。
	OK = 0

	// BadRequest 请求参数不合法（通用）。
	BadRequest = 40000

	// InvalidID 路径或参数 ID 非法。
	InvalidID = 40011

	// InvalidRequestBody 请求体 JSON 非法（通用）。
	InvalidRequestBody = 40001
	// OptionIDRequired 投票选项必填。
	OptionIDRequired = 40002
)

// UserAuth 相关业务码。
const (
	UserAuthInvalidBodySendSMS = 40041
	UserAuthPhoneRequired      = 40042
	UserAuthInvalidBodyReg     = 40043
	UserAuthPhonePwdRequired   = 40044
	UserAuthInvalidBodyLogin   = 40045
	UserAuthPhonePwdNeed       = 40046
	UserAuthInvalidBodySMS     = 40047
	UserAuthPhoneCodeRequired  = 40048
	UserAuthAccessTokenNeed    = 40049
)

// Upstream 相关网关错误码。
const (
	UpstreamUnavailable = 50201
	UpstreamEmptyReply  = 50202
	UpstreamBadPayload  = 50203
)
