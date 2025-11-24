package model

import (
	"fmt"
	"path"
	"runtime"
)

type ErrorStack string

// Error stacker for neater logging
func StackError(err error, msg string) error {
	callstr := "         ???????:??  "
	_, file, ln, ok := runtime.Caller(1) // Skip up the stack by 1 to get the calling function (not this one)
	if ok {
		filename := path.Base(file)
		callstr = fmt.Sprintf("%15s:%-4d", filename, ln)
	}

	if err == nil {
		return fmt.Errorf("%s | %s", callstr, msg)
	}

	_, is_es := err.(ErrorStack)
	if is_es {
		return ErrorStack(fmt.Sprintf("%s | %s\n%s", callstr, msg, err))
	}

	return ErrorStack(fmt.Sprintf("%s | %s\n%*s | %s", callstr, msg, len(callstr), "", err))
}

func (e ErrorStack) Error() string {
	return string(e)
}

type NotAuthenticatedError string

func (e NotAuthenticatedError) Error() string {
	return fmt.Sprintf("Can't make request; the user is not authenticated:\n%v\n", string(e))
}

type BadRequestError string

func (e BadRequestError) Error() string {
	return fmt.Sprintf("Request Failed due to incorrect formatting:\n%v\n", string(e))
}

type EmailExistsError string

func (e EmailExistsError) Error() string {
	return fmt.Sprintf("Request failed due to the provided email already being registered in the database:\n%v\n", string(e))
}

type ConnectionInvalidError string

func (e ConnectionInvalidError) Error() string {
	return fmt.Sprintf("Database connection is invalid: %s", string(e))
}

type InvalidCredentialsError string

func (e InvalidCredentialsError) Error() string {
	return fmt.Sprintf("Request failed due to the provided credentials being invalid/incorrect:\n%v\n", string(e))
}

type NotFoundError string

func (e NotFoundError) Error() string {
	return fmt.Sprintf("Couldn't find the object requested:\n%v\n", string(e))
}

type InternalServerError string

func (e InternalServerError) Error() string {
	return fmt.Sprintf("An error occurred on the server:\n%v\n", string(e))
}

type AddingItemFailed string

func (e AddingItemFailed) Error() string {
	return fmt.Sprintf("Failed to add the item given:\n%v\n", string(e))
}

type NoPasswordProvidedError int

func (e NoPasswordProvidedError) Error() string {
	return "You can't create an account without any password.\n"
}

type UnknownHttpError string

func (e UnknownHttpError) Error() string {
	return fmt.Sprintf("Some unknown non-ok status was returned by the API:\n%v\n", string(e))
}

type ForbiddenError string

func (e ForbiddenError) Error() string {
	return fmt.Sprintf("Tried to read/write something the auth-user doesn't have access to:\n%v\n", string(e))
}

type PriceOutOfRangeError uint64

func (e PriceOutOfRangeError) Error() string {
	return fmt.Sprintf("The price provided was larger than a 32-bit integer. I doubt anyone's going to buy that for you, buddy...: %v\n", uint64(e))
}

type InvalidTokenError string

func (e InvalidTokenError) Error() string {
	return fmt.Sprintf("Invalid Token format received from the API: %s\n", string(e))
}
