package parser

import (
	"errors"
	"strings"
)

var (
	// Errors
	EmptyResponseError = errors.New("Response is empty")
	EmptyTokenError    = errors.New("Token is empty")

	VariableNameNotFoundError   = errors.New("Variable name is not found")
	VariableNameNotTextError    = errors.New("Variable name cannot be another variable")
	VariableNameNotAllowedError = errors.New("Variable name is not allowed")

	VariableHasNoCloseTagError  = errors.New("Variable has no closing tag")
	VariableCanNotBeNestedError = errors.New("Variable can no be nested")
	VariableNotEnabledError     = errors.New("Variable is not enabled")
)

// Variables that can be nested in another variable as argument
var NestableVariableNames = []string{
	User, UserID, DisplayName, Channel, SubscribeLink, Commands,
}

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

/*
Validation logic:

(1) For TextTypeToken, it is always valid as long as raw text is not empty,
	because the parser does not create text token for empty string at all.
    Whitespace-only raw text is allowed.
(2) For VariableTypeToken, it should have
  [1] Valid variable name, satisfying all three conditions
	i) One of registered names
	ii) Enabled
	iii) If nested, one of nestables
  [2] Valid number of arguments (TODO)
*/
func ValidateToken(token *Token, nested bool) error {
	// For now, TextTypeToken is always valid as long as it's not empty.
	if token.TokenType == TextTokenType {
		if token.RawText == "" {
			return EmptyTokenError
		}
		return nil
	}

	// Check if variable name is not empty
	name := token.VariableName
	if name == "" {
		return VariableNameNotFoundError
	}

	// Check if variable name is allowed
	if !IsNameAllowed(name) {
		return VariableNameNotAllowedError
	}

	// If this variable is nested, check for nestability of the name
	if nested && !IsNestableName(name) {
		return VariableCanNotBeNestedError
	}

	// TODO: number of arguments check
	for _, argument := range token.Arguments {
		if err := ValidateToken(&argument, true); err != nil {
			return err
		}
	}

	// Check for ending tag
	if !strings.HasSuffix(token.RawText, VariableClosingTag) {
		return VariableHasNoCloseTagError
	}

	return nil
}

func IsNameAllowed(name string) bool {
	_, exists := VariableMap[name]
	return exists
}

func IsNestableName(name string) bool {
	// Only a few nestable variable names, just iterate through the slice
	for _, allowedName := range NestableVariableNames {
		if allowedName == name {
			return true
		}
	}
	return false
}
