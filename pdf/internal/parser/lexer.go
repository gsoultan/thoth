package parser

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"unicode"
)

type Lexer struct {
	r *bufio.Reader
}

func NewLexer(r io.Reader) *Lexer {
	return &Lexer{r: bufio.NewReader(r)}
}

func (l *Lexer) NextToken() (Token, error) {
	l.skipWhitespace()

	ch, _, err := l.r.ReadRune()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return Token{Type: TokenEOF}, nil
		}
		return Token{Type: TokenError}, err
	}

	switch ch {
	case '[':
		return Token{Type: TokenOpenBracket, Value: "["}, nil
	case ']':
		return Token{Type: TokenCloseBracket, Value: "]"}, nil
	case '/':
		return l.readName()
	case '(':
		return l.readString()
	case '<':
		next, _, _ := l.r.ReadRune()
		if next == '<' {
			return Token{Type: TokenOpenDict, Value: "<<"}, nil
		}
		l.r.UnreadRune()
		return l.readHex()
	case '>':
		next, _, _ := l.r.ReadRune()
		if next == '>' {
			return Token{Type: TokenCloseDict, Value: ">>"}, nil
		}
		l.r.UnreadRune()
		return Token{Type: TokenCloseAngle, Value: ">"}, nil
	default:
		if unicode.IsDigit(ch) || ch == '-' || ch == '+' || ch == '.' {
			l.r.UnreadRune()
			return l.readNumber()
		}
		if unicode.IsLetter(ch) {
			l.r.UnreadRune()
			return l.readKeyword()
		}
	}

	return Token{Type: TokenError, Value: string(ch)}, nil
}

func (l *Lexer) skipWhitespace() {
	for {
		ch, _, err := l.r.ReadRune()
		if err != nil {
			return
		}
		if !unicode.IsSpace(ch) && ch != 0 {
			l.r.UnreadRune()
			break
		}
		if ch == '%' {
			// Skip comment
			for {
				ch, _, err = l.r.ReadRune()
				if err != nil || ch == '\n' || ch == '\r' {
					break
				}
			}
		}
	}
}

func (l *Lexer) readName() (Token, error) {
	var buf bytes.Buffer
	for {
		ch, _, err := l.r.ReadRune()
		if err != nil {
			break
		}
		if unicode.IsSpace(ch) || isDelimiter(ch) {
			l.r.UnreadRune()
			break
		}
		buf.WriteRune(ch)
	}
	return Token{Type: TokenName, Value: buf.String()}, nil
}

func (l *Lexer) readString() (Token, error) {
	var buf bytes.Buffer
	parens := 1
	for {
		ch, _, err := l.r.ReadRune()
		if err != nil {
			break
		}
		if ch == '(' {
			parens++
		} else if ch == ')' {
			parens--
			if parens == 0 {
				break
			}
		}
		buf.WriteRune(ch)
	}
	return Token{Type: TokenString, Value: buf.String()}, nil
}

func (l *Lexer) readHex() (Token, error) {
	var buf bytes.Buffer
	for {
		ch, _, err := l.r.ReadRune()
		if err != nil || ch == '>' {
			break
		}
		// In a real PDF parser we should handle whitespace in hex strings
		if !unicode.IsSpace(ch) {
			buf.WriteRune(ch)
		}
	}
	return Token{Type: TokenHex, Value: buf.String()}, nil
}

func (l *Lexer) readNumber() (Token, error) {
	var buf bytes.Buffer
	isFloat := false
	for {
		ch, _, err := l.r.ReadRune()
		if err != nil {
			break
		}
		if !unicode.IsDigit(ch) && ch != '.' && ch != '-' && ch != '+' {
			l.r.UnreadRune()
			break
		}
		if ch == '.' {
			isFloat = true
		}
		buf.WriteRune(ch)
	}
	val := buf.String()
	if isFloat {
		return Token{Type: TokenFloat, Value: val}, nil
	}
	return Token{Type: TokenInteger, Value: val}, nil
}

func (l *Lexer) readKeyword() (Token, error) {
	var buf bytes.Buffer
	for {
		ch, _, err := l.r.ReadRune()
		if err != nil {
			break
		}
		if !unicode.IsLetter(ch) {
			l.r.UnreadRune()
			break
		}
		buf.WriteRune(ch)
	}
	return Token{Type: TokenKeyword, Value: buf.String()}, nil
}

func (l *Lexer) ReadRaw(length int) ([]byte, error) {
	// Skip the newline following the 'stream' keyword.
	// Common formats: \n or \r\n
	ch, _, err := l.r.ReadRune()
	if err != nil {
		return nil, err
	}
	if ch == '\r' {
		next, _, _ := l.r.ReadRune()
		if next != '\n' {
			l.r.UnreadRune()
		}
	} else if ch != '\n' {
		// If it's not a newline, we should probably unread it,
		// though standard PDF says it must be a newline.
		l.r.UnreadRune()
	}

	data := make([]byte, length)
	_, err = io.ReadFull(l.r, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func isDelimiter(ch rune) bool {
	switch ch {
	case '(', ')', '<', '>', '[', ']', '{', '}', '/', '%':
		return true
	}
	return false
}
