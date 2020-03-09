package models

const (
	// TextTokenType pure text token in parsed response
	TextTokenType = 1
	// VariableTokenType variable token, in form of $(name [args0, args1, ...])
	// For example, $(channel_id), $(countup OOO), etc. are of variable token type
	VariableTokenType = 2
)

// ParsedResponse response raw text is parsed into tokens
// TODO: Any better name than "parsed reponse"?
type ParsedResponse struct {
	RawText string // RawText is simply concated RawText of tokens. Is it really needed?
	Tokens  []Token
}

/*
Token a response is parsed into a sequence of tokens, which is either a text or a variable.
Each token is then passed into converter, which resolves variables and completes the response text

For example, "Welcome, $(user) to the stream. Stream is up for $(uptime)"
is parsed to three tokens {"Welcome, ", "$(user)", " to the stream. Stream is up for ", "$(uptime)"}
Note that

Tokens "Welcome, " and " to the stream. Stream is up for " are of TextToken type
Tokens "$(user)" and "$(uptime)" are of VariableToken type

For TextToken objects, token.RawText is the value.
For VariableToken objects, the value should be somehow retrieved from somewhere.

The type of arguments is Token, not string, because of nested variables
*/
type Token struct {
	RawText      string  // raw text of the token, including arguments and nested variables.
	TokenType    int     // string or variable.
	VariableName string  // needed to find Variable struct with this name.
	Arguments    []Token // for function variables like $(time [arg1]), $(rand [arg1], [arg2]), etc.
}
