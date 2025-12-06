package errcode

import "fmt"

// Error 自定义错误类型
type Error struct {
	Code int
	Msg  string
}

func (e *Error) Error() string {
	return e.Msg
}

func NewError(code int, msg string) *Error {
	return &Error{Code: code, Msg: msg}
}

// 通用错误 (1xxxx)
var (
	ErrSuccess         = NewError(0, "success")
	ErrInvalidParam    = NewError(10001, "参数错误")
	ErrInternalServer  = NewError(10002, "服务器内部错误")
	ErrUnauthorized    = NewError(10003, "未授权")
	ErrForbidden       = NewError(10004, "无权限")
	ErrNotFound        = NewError(10005, "资源不存在")
	ErrTooManyRequests = NewError(10006, "请求过于频繁")
)

// 用户相关 (2xxxx)
var (
	ErrUserNotFound       = NewError(20001, "用户不存在")
	ErrUserExists         = NewError(20002, "用户已存在")
	ErrInvalidCode        = NewError(20003, "验证码错误")
	ErrCodeExpired        = NewError(20004, "验证码已过期")
	ErrInvalidPhone       = NewError(20005, "手机号格式错误")
	ErrSendCodeTooFreq    = NewError(20006, "发送验证码过于频繁")
	ErrInvalidToken       = NewError(20007, "Token 无效")
	ErrTokenExpired       = NewError(20008, "Token 已过期")
	ErrSubscriptionExpire = NewError(20009, "订阅已过期")
)

// 分身相关 (3xxxx)
var (
	ErrAvatarNotFound     = NewError(30001, "分身不存在")
	ErrAvatarExists       = NewError(30002, "分身已存在")
	ErrAvatarNameInvalid  = NewError(30003, "分身名称不合法")
	ErrAvatarLimitReached = NewError(30004, "分身数量已��上限")
	ErrPersonalityInvalid = NewError(30005, "性格参数不合法")
)

// AI 相关 (4xxxx)
var (
	ErrAIServiceFailed  = NewError(40001, "AI 服务调用失败")
	ErrPromptNotFound   = NewError(40002, "Prompt 模板不存在")
	ErrAITimeout        = NewError(40003, "AI 服务超时")
	ErrAIQuotaExceeded  = NewError(40004, "AI 调用额度已用完")
	ErrInvalidPrompt    = NewError(40005, "Prompt 参数不合法")
	ErrEmbeddingFailed  = NewError(40006, "文本向量化失败")
)

// 事件相关 (5xxxx)
var (
	ErrEventNotFound       = NewError(50001, "事件不存在")
	ErrEventTemplateNotFound = NewError(50002, "事件模板不存在")
	ErrEventGenerateFailed = NewError(50003, "事件生成失败")
)

// 世界相关 (6xxxx)
var (
	ErrRegionNotFound = NewError(60001, "区域不存在")
	ErrSceneNotFound  = NewError(60002, "场景不存在")
	ErrSceneLocked    = NewError(60003, "场景未解锁")
)

// 对话相关 (7xxxx)
var (
	ErrDialogueNotFound = NewError(70001, "对话不存在")
	ErrMessageTooLong   = NewError(70002, "消息过长")
	ErrMessageEmpty     = NewError(70003, "消息不能为空")
)

// 日记相关 (8xxxx)
var (
	ErrDiaryNotFound      = NewError(80001, "日记不存在")
	ErrDiaryExists        = NewError(80002, "今日日记已存在")
	ErrDiaryContentEmpty  = NewError(80003, "日记内容不能为空")
	ErrDiaryGenerateFailed = NewError(80004, "日记生成失败")
)

// OSS 相关 (9xxxx)
var (
	ErrOSSUploadFailed   = NewError(90001, "文件上传失败")
	ErrOSSDeleteFailed   = NewError(90002, "文件删除失败")
	ErrFileTypeInvalid   = NewError(90003, "文件类型不支持")
	ErrFileSizeExceeded  = NewError(90004, "文件大小超出限制")
	ErrOSSTokenFailed    = NewError(90005, "获取上传凭证失败")
)

// 短信相关 (11xxxx)
var (
	ErrSMSSendFailed = NewError(110001, "短信发送失败")
	ErrSMSVerifyFailed = NewError(110002, "短信验证失败")
)

// 调度相关 (12xxxx)
var (
	ErrScheduleNotFound = NewError(120001, "调度配置不存在")
	ErrScheduleFailed   = NewError(120002, "调度执行失败")
	ErrSchedulePaused   = NewError(120003, "调度已暂停")
)

// FromError 从 error 转换为 Error
func FromError(err error) *Error {
	if err == nil {
		return ErrSuccess
	}
	if e, ok := err.(*Error); ok {
		return e
	}
	return NewError(10002, err.Error())
}

// Wrap 包装错误信息
func Wrap(err *Error, msg string) *Error {
	return NewError(err.Code, fmt.Sprintf("%s: %s", err.Msg, msg))
}
