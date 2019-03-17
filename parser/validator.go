package parser

import (
	"errors"
	"strings"
)

var (
	EmptyResponseError          = errors.New("Response is empty")
	VariableNameNotFoundError   = errors.New("Variable name is not found")
	VariableNameNotTextError    = errors.New("Variable name cannot be another variable")
	VariableNameNotAllowedError = errors.New("Variable name is not allowed")
	VariableHasNoCloseTagError  = errors.New("Variable has no closing tag")
	VariableCanNotBeNestedError = errors.New("Variable can no be nested")
)

/*
Validation logic:

(1) For TextTypeToken, it is always valid as long as it is not empty string token
(2) For VariableTypeToken, it should have
  [1] Valid variable name (nested or not nested)
  [2] Valid number of arguments
*/
// Returns nil if there is no error
func Validate(response *ParsedResponse) error {
	if response.RawText == "" {
		return EmptyResponseError
	}

	for _, token := range response.Tokens {
		if err := ValidateToken(&token, false); err != nil {
			return err
		}
	}
	return nil
}

// TODO: Variable type with no token is valid?
// For example, "$(user)" has no token in it.
func ValidateToken(token *Token, nested bool) error {
	// For now, TextTypeToken is always valid.
	if token.TokenType == TextTokenType {
		return nil
	}

	// Check if variable name is correct
	name, err := GetVariableName(token)
	if err != nil {
		return err
	}

	// Check if variable name is allowed
	if !IsNameAllowed(name) {
		return VariableNameNotAllowedError
	}

	// Check if this variable can be nested in another variable
	if nested && !NameCanBeNested(name) {
		return VariableCanNotBeNestedError
	}

	for _, argument := range token.Arguments {
		if err := ValidateToken(&argument, true); err != nil {
			return err
		}
	}

	// Check for ending tag
	if !strings.HasSuffix(token.RawText, VariableCloseTag) {
		return VariableHasNoCloseTagError
	}

	return nil
}

// TODO: Eventually, it should be one of allowed variable names.
func IsNameAllowed(name string) bool {
	return true
}

// TODO: Implement this
func NameCanBeNested(name string) bool {
	return true
}

// GetVariableName returns
func GetVariableName(token *Token) (string, error) {
	// Variable token with empty body "$()"
	if len(token.Arguments) == 0 {
		return "", VariableNameNotFoundError
	}
	nameArg := token.Arguments[0]
	// The first argument of the variable token should always be text.
	// TokenType not being TextTokenType means nested variable in form of $($(a) b)
	if nameArg.TokenType != TextTokenType {
		return "", VariableNameNotTextError
	}

	fields := strings.Fields(nameArg.RawText)
	if len(fields) == 0 {
		return "", VariableNameNotFoundError
	}
	return fields[0], nil
}
