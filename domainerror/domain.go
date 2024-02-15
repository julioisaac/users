package domainerror

// ErrorCode
type ErrorCode string

const (
	ErrorCodeNone        ErrorCode = ""
	Default              ErrorCode = "USR001"
	ValidationError      ErrorCode = "USR002"
	InvalidRequestData   ErrorCode = "USR003"
	CreateUserError      ErrorCode = "USR004"
	GetUserDefaultError  ErrorCode = "USR005"
	GetUserNotFoundError ErrorCode = "USR006"
)

// DomainError
type Error struct {
	Code    ErrorCode
	Message string
	Detail  map[string]any
}

func (m *Error) Error() string {
	return m.Message
}

// New domain error
func New(code ErrorCode, message string, detail map[string]any) *Error {
	err := new(Error)

	err.Code = code
	err.Message = message
	err.Detail = detail

	return err
}
