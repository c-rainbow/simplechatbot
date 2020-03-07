package chatplugins

import "errors"

const (
	// DefaultFailureMessage default error message
	DefaultFailureMessage = "Failed to run command"
)

// Some common error types
var (
	// ErrCommandNotFound error when the command name is not found
	ErrCommandNotFound = errors.New("Command with the name is not found")
	// ErrNoPermission error when the chatter does not have permission to use this command
	ErrNoPermission = errors.New("User has no permission")
	// ErrNotEnoughArguments error when the chatter didn't provide enough arguments for command
	ErrNotEnoughArguments = errors.New("Not enough arguments")

	// ErrTargetCommandAlreadyExists used by AddCommand
	ErrTargetCommandAlreadyExists = errors.New("Target command already exists")
	// ErrTargetCommandNotFound used by EditCommand, DeleteCommand
	ErrTargetCommandNotFound = errors.New("Target command does not exist")

	// ErrInvalidArgument one or more provided arguments are invalid
	ErrInvalidArgument = errors.New("Invalid argument")

	// ErrResponseNotFound no response is found for the response key
	ErrResponseNotFound = errors.New("ResponseKey has no matching response object")
)

// KeyedErrorT error with string key code
type KeyedErrorT interface {
	Key() string
	Error() string
}

// KeyedError error with string key code
type KeyedError struct {
	key     string
	message string
}

var _ KeyedErrorT = (*KeyedError)(nil)

// Key returns string error key
func (err *KeyedError) Key() string {
	return err.key
}

// Error returns string error message
func (err *KeyedError) Error() string {
	return err.message
}

// NewError creates a new keyed error
func NewError(key string, message string) KeyedErrorT {
	return &KeyedError{key: key, message: message}
}
