package parser

type TokenType int

const (
	TokenError TokenType = iota
	TokenEOF
	TokenKeyword
	TokenInteger
	TokenFloat
	TokenString
	TokenHex
	TokenName
	TokenOpenBracket  // [
	TokenCloseBracket // ]
	TokenOpenBrace    // {
	TokenCloseBrace   // }
	TokenOpenDict     // <<
	TokenCloseDict    // >>
	TokenCloseAngle   // >
)
