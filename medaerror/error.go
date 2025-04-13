package medaerror

import "fmt"

const (
	STANDARD_ERROR = 1313
	NO_ERROR       = 8999 // PLEASE MAKE SURE THIS IS NOT USED ELSE WHERE!!!
)

// Some function aliases to make it simpler and backward compatible
var (
	New       func(int, string, string, interface{}) *MedaError = NewMedaErrPtr
	Errorf    func(string, ...any) MedaError                    = NewMedaErrStringFormat // error that only contains the error string with format, no code, no data
	NewString func(string) MedaError                            = NewMedaErrString       // error that only contains the error string, no code, no data
	Simple    func(string) MedaError                            = NewMedaErrString
)

// This is wrapper for go standard error
type MedaError struct {
	Code           int
	Message        string // usually for logging
	ResponseString string
	Data           interface{}
	FunctionName   string
	IsLogged       bool
	Err            error // original error, if this is nil, then MedaError is used for logging only
}

func NewMedaErr(code int, message, respString string, data interface{}) MedaError {
	if respString == "" {
		respString = message
	}
	return MedaError{
		Code:           code,
		Message:        message,
		ResponseString: respString,
		Data:           data,
	}
}

func NewMedaErrString(message string) MedaError {
	return NewMedaErr(STANDARD_ERROR, message, message, nil)
}

func NewMedaErrStringFormat(format string, a ...any) MedaError {
	msg := fmt.Sprintf(format, a...)
	return NewMedaErr(STANDARD_ERROR, msg, msg, nil)
}

// Same with NewMedaErr but return pointer
// NOTE: do we need this?
func NewMedaErrPtr(code int, message, respString string, data interface{}) *MedaError {
	err := NewMedaErr(code, message, respString, data)
	return &err
}

// Same with NewMedaErr but return pointer
// NOTE: do we need this?
func NewMedaErrPtrString(message string) *MedaError {
	err := Errorf(message)
	return &err
}

// To implement standard golang Error type
func (m MedaError) Error() string {
	return m.Message
}

// Return the response string usually that needed for API handler
func (m MedaError) Response() string {
	return m.ResponseString
}

// Yes this is the same as Error() but for readibility sometimes use Log() for logging is clearer
func (m MedaError) Log() string {
	return m.Message
}

// Print both the mssage and response
func (m MedaError) String() string {
	return fmt.Sprintf("code:%d;message:%s;response:%s;data:%v\n", m.Code, m.Message, m.ResponseString, m.Data)
}

// Print both message and response in pretty way with newlines, possibly for debugging
func (m MedaError) Pretty() string {
	return fmt.Sprintf("Code:%d\nMessage:%s\nResponse:%s\nData:%v\n", m.Code, m.Message, m.ResponseString, m.Data)
}
