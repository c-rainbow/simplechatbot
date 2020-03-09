package parser

import (
	"errors"
	"strconv"
	"strings"
	"time"

	l10n "github.com/c-rainbow/simplechatbot/localization/common"
	models "github.com/c-rainbow/simplechatbot/models"
	"github.com/c-rainbow/simplechatbot/resolver"
	twitch_irc "github.com/gempir/go-twitch-irc"
)

var (
	NotVariableTypeError         = errors.New("Not variable type token")
	InvalidVariableNameError     = errors.New("Invalid variable name")
	UnsupportedVariableTypeError = errors.New("Unsupported variable type")
)

// It is assumed that the response is already validated
func ConvertResponse(
	response *models.ParsedResponse, channel string, sender *twitch_irc.User, message *twitch_irc.PrivateMessage,
	locale *l10n.LocaleConfig, args []string) (string, error) {
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
				converted, err = ConvertStreamAPIVariables(&token, channel, sender, message, locale, args)
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
	token *models.Token, channel string, sender *twitch_irc.User, message *twitch_irc.PrivateMessage) (string, error) {
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
		return sender.ID, nil
	case DisplayName:
		return sender.DisplayName, nil
	case Channel:
		return channel, nil
	// TODO: Don't hardcode link
	case SubscribeLink:
		return "https://twitch.tv/subs/" + channel, nil
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

func ConvertStreamAPIVariables(token *models.Token, channel string, sender *twitch_irc.User,
	message *twitch_irc.PrivateMessage, locale *l10n.LocaleConfig, args []string) (string, error) {
	if token.TokenType != models.VariableTokenType {
		return "", NotVariableTypeError
	}

	resolver := resolver.DefaultStreamsAPIResolver()

	stream, err := resolver.Resolve(channel)
	if err != nil {
		return "", err
	}

	switch token.VariableName {
	// This case should be updated if more argument indexes are added
	case Title:
		return stream.Title, nil
	case Game:
		// This is Game name or Game ID?
		return stream.GameID, nil
	case Uptime:
		uptimeDuration := time.Since(stream.StartedAt)
		return locale.DurationToString(uptimeDuration), nil
	case UptimeAt:
		startTime := stream.StartedAt
		return locale.DateTimeToString(startTime), nil
	case ViewerCount:
		return strconv.Itoa(stream.ViewerCount), nil
	case SubscriberCount:
		return "", VariableNotEnabledError
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
