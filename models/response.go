package models

const (
	TextTokenType     = 1
	VariableTokenType = 2
)

// TODO: Any better name than "parsed reponse"?
type ParsedResponse struct {
	RawText string // RawText is simply concated RawText of tokens. Is it really needed?
	Tokens  []Token
}

/*
For example, "Welcome, $(user) to the stream. Stream is up for $(uptime)"
is parsed to three tokens {"Welcome, ", "$(user)", " to the stream. Stream is up for ", "$(uptime)"}
Note that

Tokens "Welcome, " and " to the stream. Stream is up for " are of TextToken type
Tokens "$(user)" and "$(uptime)" are of VariableToken type

For TextToken objects, token.RawText is the value.
For VariableToken objects, the value should be somehow retrieved from somewhere.

The type of arguments is Token, not string, because of nested variables

TODO: Create another struct called NestedToken, or re-use the same Token struct?
*/
type Token struct {
	RawText      string  // raw text of the token, including arguments and nested variables.
	TokenType    int     // string or variable.
	VariableName string  // needed to find Variable struct with this name.
	Arguments    []Token // for function variables like $(time [arg1]), $(rand [arg1], [arg2]), etc.
}
