package wmenu

import "errors"

var (
	//ErrInvalid is returned if a response from user was an invalid option
	ErrInvalid = errors.New("invalid response")

	//ErrTooMany is returned if multiSelect is false and a user tries to select multiple options
	ErrTooMany = errors.New("too many responses")

	//ErrNoResponse is returned if there were no responses and no action to call
	ErrNoResponse = errors.New("no response")

	//ErrDuplicate is returned is a user selects an option twice
	ErrDuplicate = errors.New("duplicated response")
)

//MenuError records menu errors
type MenuError struct {
	Err       error
	Res       string
	TriesLeft int
}

//Error prints the error in an easy to read string.
func (e *MenuError) Error() string {
	if e.Res != "" {
		return e.Err.Error() + ": " + e.Res
	}
	return e.Err.Error()
}

func newMenuError(err error, res string, tries int) *MenuError {
	return &MenuError{
		Err:       err,
		Res:       res,
		TriesLeft: tries,
	}
}

//IsInvalidErr checks to see if err is of type invalid error returned by menu.
func IsInvalidErr(err error) bool {
	e, ok := err.(*MenuError)
	if ok && e.Err == ErrInvalid {
		return true
	}
	return false
}

//IsNoResponseErr checks to see if err is of type no response returned by menu.
func IsNoResponseErr(err error) bool {
	e, ok := err.(*MenuError)
	if ok && e.Err == ErrNoResponse {
		return true
	}
	return false
}

//IsTooManyErr checks to see if err is of type too many returned by menu.
func IsTooManyErr(err error) bool {
	e, ok := err.(*MenuError)
	if ok && e.Err == ErrTooMany {
		return true
	}
	return false
}

//IsDuplicateErr checks to see if err is of type duplicate returned by menu.
func IsDuplicateErr(err error) bool {
	e, ok := err.(*MenuError)
	if ok && e.Err == ErrDuplicate {
		return true
	}
	return false
}

//IsMenuErr checks to see if it is a menu err.
//This is a general check not a specific one.
func IsMenuErr(err error) bool {
	_, ok := err.(*MenuError)
	if ok {
		return true
	}
	return false
}
