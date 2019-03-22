package chatplugins

import "errors"

const (
	DefaultFailureMessage = "Failed to run command"
)

var (
	ErrCommandNotFound    = errors.New("Command with the name is not found")
	ErrNoPermission       = errors.New("User has no permission")
	ErrNotEnoughArguments = errors.New("Not enough arguments")

	ErrTargetCommandAlreadyExists = errors.New("Target command already exists")
	ErrTargetCommandNotFound      = errors.New("Target command does not exist")

	ErrInvalidArgument = errors.New("Invalid argument")

	ErrResponseNotFound = errors.New("ResponseKey has no matching response object")
)

type KeyedErrorT interface {
	Key() string
	Error() string
}

type KeyedError struct {
	key     string
	message string
}

var _ KeyedErrorT = (*KeyedError)(nil)

func (err *KeyedError) Key() string {
	return err.key
}

func (err *KeyedError) Error() string {
	return err.message
}

func NewError(key string, message string) KeyedErrorT {
	return &KeyedError{key: key, message: message}
}
