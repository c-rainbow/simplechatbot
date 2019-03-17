package parser

import (
	"errors"
	"strings"

	twitch_irc "github.com/gempir/go-twitch-irc"
)

var (
	NotVariableTypeError         = errors.New("Not variable type token")
	InvalidVariableNameError     = errors.New("Invalid variable name")
	UnsupportedVariableTypeError = errors.New("Unsupported variable type")
)

// It is assumed that the response is already validated
func ConvertResponse(response *ParsedResponse, channel string, sender *twitch_irc.User, message *twitch_irc.Message) (string, error) {
	var builder strings.Builder
	for _, token := range response.Tokens {
		if token.TokenType == TextTokenType {
			builder.WriteString(token.RawText)
		} else {
			var converted string
			var err error

			// There seems to be no other way than manually doing variable type check
			variableType := VariableMap[token.VariableName].Type
			switch variableType {
			case ChatType:
				converted, err = ConvertChatVariables(&token, channel, sender, message)

			case StreamAPIType:
				fallthrough
			case UserAPIType:
				fallthrough
			case SimpleType:
				fallthrough
			case OverwatchAPIType:
				fallthrough
			case LeagueOfLegendsAPIType:
				// TODO: Change this to an appropriate type
				return "", VariableNotEnabledError
			default:
				return "", UnsupportedVariableTypeError
			}
			if err == nil {
				builder.WriteString(converted)
			} else {
				return "", err
			}
		}
	}
	return builder.String(), nil
}

func ConvertChatVariables(token *Token, channel string, sender *twitch_irc.User, message *twitch_irc.Message) (string, error) {
	// This check is here only for robustness.
	// In normal workflow, this if-statement will always be skipped.
	if token.TokenType != VariableTokenType {
		return "", NotVariableTypeError
	}

	// Unfortunately, this giant switch seems to be the only way
	switch token.VariableName {
	case User:
		return sender.DisplayName, nil
	case UserID:
		return sender.UserID, nil
	case DisplayName:
		return sender.DisplayName, nil
	case Channel:
		return channel, nil
	// TODO: Don't hardcode link
	case SubscribeLink:
		return "http://twitch.tv/subs/" + channel, nil
	default:
		return "", InvalidVariableNameError
	}
}
