package golang_social_chat

import (
	"encoding/json"
	"path"
	"runtime"
)

// _BaseError defines the basic information that an error carries.
type _BaseError struct {
	Status  uint           `json:"status"`
	Message string         `json:"message"`
	Data    map[string]any `json:"data"`
}

// Error answers the message. For satisfying the `error` interface.
func (e *_BaseError) Error() string {
	buf, _ := json.Marshal(e)
	return string(buf)
}

// Unwrap answers the wrapped error.
func (e *_BaseError) Unwrap() error {
	if e.Data != nil {
		if _err, ok := e.Data["_inner"]; ok {
			return _err.(error)
		}
	}
	return nil
}

// UserError represents an error that is caused by some input condition.
type UserError struct {
	_BaseError
}

// SystemError represents an error that is *not* caused by some input condition.
// but is detected deep down.
type SystemError struct {
	_BaseError
}

// NewUserError creates a new error that represents an input-driven error condition.
func NewUserError(status uint, msg string, err error) *UserError {
	if status == 0 || (msg == "" && err == nil) {
		panic("invalid UserError creation attempt")
	}

	_err := &UserError{}
	_err.Status = status
	_err.Message = msg
	if err != nil {
		_err.Data = map[string]any{}
		switch err.(type) {
		case *UserError, *SystemError:
			_err.Data["_inner"] = err
		default:
			_err.Data["_inner"] = err.Error()
		}
	}
	return _err
}

// Add sets a new key-value error data pair.
func (e *UserError) Add(key string, value any) *UserError {
	if e.Data == nil {
		e.Data = map[string]any{}
	}
	e.Data[key] = value
	return e
}

// NewSystemError create a new error that represents a system error condition.
func NewSystemError(status uint, msg string, err error) *SystemError {
	if status == 0 || (msg == "" && err == nil) {
		panic("invalid SystemError creation attempt")
	}

	_err := &SystemError{}
	_err.Status = status
	_err.Message = msg
	_err.Data = map[string]any{}
	if err != nil {
		switch err.(type) {
		case *UserError, *SystemError:
			_err.Data["_inner"] = err
		default:
			_err.Data["_inner"] = err.Error()
		}
	}
	if _, _file, _line, ok := runtime.Caller(1); ok {
		_err.Data["_file"] = path.Base(_file)
		_err.Data["_line"] = _line
	}
	return _err
}

// Add set a new key-value error data pair.
func (e *SystemError) Add(key string, value any) *SystemError {
	e.Data[key] = value
	return e
}

// Result holds the results of a successful processing of a user request.
type Result struct {
	status  uint
	message string
	data    map[string]any
}

// NewResult creates a result that holds the success result of an operation.
func NewResult(status uint, msg string) *Result {
	if status == 0 || msg == "" {
		panic("invalid Result creation attempt")
	}
	_res := &Result{}
	_res.status = status
	_res.message = msg
	return _res
}

// Add sets a new key-vale data pair.
func (r *Result) Add(key string, value any) *Result {
	if r.data == nil {
		r.data = map[string]any{}
	}
	r.data[key] = value
	return r
}
