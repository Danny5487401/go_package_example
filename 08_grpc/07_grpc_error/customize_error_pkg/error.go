package customize_error_pkg

import "net/http"

// 这里可以是公司定义的公用错误
var (
	Err44010102 = NewError(44010102, "请求参数错误")
	Err44010101 = NewError(44010101, "网络异常")
)

var _ ErrorNo = (*err)(nil)

type ErrorNo interface {
	i()
	WithTitle(title string) ErrorNo
	GetTitle() string
	GetCode() int
	Error() string
}
type err struct {
	Errors errors `json:"errors"`
}

//title,details 两个字段的内容一样，现在的客户端正常只用其中一个，为了兼容老版本客户端，此处继续维护
type errors struct {
	Id         int    `json:"id,string"`     // id
	Code       int    `json:"code,string"`   // 业务编码
	Level      Level  `json:"level,string"`  // 弹窗提示类型  详解 constant文件
	Status     int    `json:"status,string"` // 状态 默认是http.StatusOK 即200
	Title      string `json:"title"`         // level= toast标题
	PopupTitle string `json:"popup_title"`   // level=2 弹框标题
	Details    string `json:"details"`       // level=2 弹框内容
}

//level=1  toast提示
func NewError(code int, title string, opts ...Option) ErrorNo {
	if len(title) == 0 {
		title = "网络异常"
	}
	e := &err{Errors: errors{
		Code:    code,
		Id:      0,
		Level:   LevelToast,
		Status:  http.StatusOK,
		Title:   title,
		Details: title,
	}}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

//level=2  弹框提示
func NewPopupError(code int, title string, detail string, opts ...Option) ErrorNo {
	if len(title) == 0 {
		title = "网络异常"
	}
	e := &err{Errors: errors{
		Code:       code,
		Id:         0,
		Level:      LevelPopup,
		Status:     http.StatusOK,
		Title:      detail, //为了兼容老版本，所以这个字段还是要返回，但是返回的标题是details
		PopupTitle: title,
		Details:    detail,
	}}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

func (e *err) i() {}

func (e *err) WithTitle(msg string) ErrorNo {
	eCopy := &err{Errors: errors{
		Code:    e.Errors.Code,
		Id:      e.Errors.Id,
		Details: e.Errors.Details,
		Level:   e.Errors.Level,
		Status:  e.Errors.Status,
		Title:   e.Errors.Title,
	}}
	eCopy.Errors.Title = msg

	return eCopy
}
func (e *err) GetTitle() string {
	return e.Errors.Title
}
func (e *err) GetCode() int {
	return e.Errors.Code
}

//仅兼容官方包error签名
func (e *err) Error() string {
	return ""
}

type Option func(e *err)

func WithLevel(level Level) Option {
	return func(e *err) {
		e.Errors.Level = level
	}
}

func WithDetail(detail string) Option {
	return func(e *err) {
		e.Errors.Details = detail
	}
}

type Level int

//优先推荐使用level 1-仅toast提示
const (
	LevelToast Level = 1 //toast提示
	LevelPopup Level = 2 //弹窗提示
)
