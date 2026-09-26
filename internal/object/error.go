package object

const KindError Kind = "错误"

type Error struct {
	Message string
}

func (e *Error) Kind() Kind      { return KindError }
func (e *Error) Inspect() string { return "错误：" + e.Message }
func (e *Error) Error() string   { return e.Message }

func NewError(message string) *Error {
	return &Error{Message: message}
}
