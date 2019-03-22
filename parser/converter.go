package parser

import (
	"errors"
	"strconv"
	"strings"

	models "github.com/c-rainbow/simplechatbot/models"
	twitch_irc "github.com/gempir/go-twitch-irc"
)

var (
	NotVariableTypeError         = errors.New("Not variable type token")
	InvalidVariableNameError     = errors.New("Invalid variable name")
	UnsupportedVariableTypeError = errors.New("Unsupported variable type")
)

// It is assumed that the response is already validated
func ConvertResponse(
	response *models.ParsedResponse, channel string, sender *twitch_irc.User, message *twitch_irc.Message,
	args []string) (string, error) {
	var builder strings.Builder
	for _, token := range response.Tokens {
		if token.TokenType == models.TextTokenType {
			builder.WriteString(token.RawText)
		} else {
			var converted string
			var err error

			// There seems to be no other way than manually doing variable type check
			variableType := VariableMap[token.VariableName].Type
			switch variableType {
			case ChatType:
				converted, err = ConvertChatVariables(&token, channel, sender, message)
			case ArgumentType:
				converted, err = ConvertArgumentVariables(&token, args)
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

func ConvertChatVariables(
	token *models.Token, channel string, sender *twitch_irc.User, message *twitch_irc.Message) (string, error) {
	// This check is here only for robustness.
	// In normal workflow, this if-statement will always be skipped.
	if token.TokenType != models.VariableTokenType {
		return "", NotVariableTypeError
	}

	// Unfortunately, this giant switch-statement seems to be the only way
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

func ConvertArgumentVariables(token *models.Token, args []string) (string, error) {
	if token.TokenType != models.VariableTokenType {
		return "", NotVariableTypeError
	}

	switch token.VariableName {
	// This case should be updated if more argument indexes are added
	case Arg0, Arg1, Arg2, Arg3, Arg4, Arg5:
		return valueAtIndex(args, token.VariableName)
	default:
		return "", InvalidVariableNameError
	}
}

// This function assumes that the variable name is form of "argN", where N >= 0
func valueAtIndex(args []string, vName string) (string, error) {
	index, err := strconv.Atoi(vName[3:])
	if err == nil && 0 <= index && index < len(args) {
		return args[index], nil
	}
	return "", InvalidVariableNameError
}
