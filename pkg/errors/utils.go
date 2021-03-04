package errors

import (
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// MySQLErrorConversion ...
func MySQLErrorConversion(err error) error {
	if err == nil {
		return nil
	}

	if err == gorm.ErrRecordNotFound {
		return ErrResourceNotFound
	}

	mysqlErr, ok := err.(*mysql.MySQLError)
	if ok {
		if mysqlErr.Number == 1062 {
			// the duplicate key error.
			return ErrResourceAlreadyExists
		}
	}

	return ErrInternalServerError
}

// HTTPError represents an error that occurred while handling a request.
type HTTPError struct {
	Code     string      `json:"code"`
	Message  interface{} `json:"message"`
	HTTPCode int         `json:"-"`
}

// GetHTTPError get http error
func GetHTTPError(err error) HTTPError {
	e, ok := err.(*_err)
	if !ok {
		return HTTPError{
			Code:     ErrInternalServerError.Code,
			Message:  ErrInternalServerError.Message,
			HTTPCode: ErrInternalServerError.HTTPCode,
		}
	}
	return HTTPError{
		Code:     e.Code,
		Message:  e.Message,
		HTTPCode: e.HTTPCode,
	}
}
