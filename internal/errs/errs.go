package errs

import (
	"fmt"
	"net/http"

	"github.com/HardDie/event_tracker/internal/logger"
)

var (
	InternalError  = NewError("internal error")
	BadRequest     = NewError("bad request", http.StatusBadRequest)
	UserBlocked    = NewError("user is blocked", http.StatusUnauthorized)
	SessionInvalid = NewError("session invalid", http.StatusUnauthorized)
)

type Err struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Err     error  `json:"err"`
}

func NewError(message string, code ...int) *Err {
	err := &Err{
		Message: message,
		Code:    http.StatusInternalServerError,
	}
	if len(code) > 0 {
		err.Code = code[0]
	}
	return err
}

func (e Err) Error() string {
	return fmt.Sprintf("HTTP[%d] %s", e.GetCode(), e.GetMessage())
}
func (e Err) Unwrap() error {
	return e.Err
}

func (e *Err) HTTP(code int) *Err {
	return &Err{
		Message: e.Message,
		Code:    code,
		Err:     e,
	}
}
func (e *Err) AddMessage(message string) *Err {
	return &Err{
		Message: message,
		Code:    e.Code,
		Err:     e,
	}
}

func (e *Err) GetCode() int       { return e.Code }
func (e *Err) GetMessage() string { return e.Message }

func HttpError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}
	val, ok := err.(*Err)
	if !ok || val == nil {
		logger.Error.Println("Unknown error:", err.Error())
		http.Error(w, "Unknown error", http.StatusInternalServerError)
		return
	}
	http.Error(w, val.Message, val.Code)
}
