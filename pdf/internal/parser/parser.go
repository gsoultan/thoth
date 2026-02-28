package parser

import (
	"fmt"
	"io"
	"strconv"

	"github.com/gsoultan/thoth/pdf/internal/objects"
)

type Parser struct {
	lexer *Lexer
	peeks []Token
}

func NewParser(lexer *Lexer) *Parser {
	return &Parser{lexer: lexer, peeks: make([]Token, 0, 2)}
}

func (p *Parser) ParseObject() (objects.Object, error) {
	tok, err := p.nextToken()
	if err != nil {
		return nil, err
	}

	switch tok.Type {
	case TokenEOF:
		return nil, io.EOF
	case TokenInteger:
		// Check if it's a reference (n m R) or indirect object (n m obj)
		peek1, err := p.peekTokenN(1)
		if err == nil && peek1.Type == TokenInteger {
			peek2, err := p.peekTokenN(2)
			if err == nil {
				if peek2.Type == TokenKeyword && (peek2.Value == "R" || peek2.Value == "obj") {
					p.nextToken() // skip peek1 (m)
					p.nextToken() // skip peek2 (R or obj)
					n, _ := strconv.Atoi(tok.Value)
					m, _ := strconv.Atoi(peek1.Value)
					if peek2.Value == "R" {
						return objects.Reference{Number: n, Generation: m}, nil
					}
					data, err := p.ParseObject()
					if err != nil {
						return nil, err
					}
					p.skipEndObj()
					return &objects.IndirectObject{Number: n, Generation: m, Data: data}, nil
				}
			}
		}
		val, _ := strconv.Atoi(tok.Value)
		return objects.Integer(val), nil
	case TokenFloat:
		val, _ := strconv.ParseFloat(tok.Value, 64)
		return objects.Float(val), nil
	case TokenName:
		return objects.Name(tok.Value), nil
	case TokenString:
		return objects.PDFString(tok.Value), nil
	case TokenOpenBracket:
		return p.parseArray()
	case TokenOpenDict:
		dict, err := p.parseDictionary()
		if err != nil {
			return nil, err
		}

		// Check if it's a stream (dictionary followed by 'stream' keyword)
		peek, err := p.peekTokenN(1)
		if err == nil && peek.Type == TokenKeyword && peek.Value == "stream" {
			p.nextToken() // consume 'stream'
			length := 0
			if lenObj, ok := dict["Length"]; ok {
				if i, ok := lenObj.(objects.Integer); ok {
					length = int(i)
				}
			}
			data, err := p.lexer.ReadRaw(length)
			if err != nil {
				return nil, err
			}
			// Skip 'endstream'
			p.skipKeyword("endstream")
			return objects.Stream{Dict: dict, Data: data}, nil
		}
		return dict, nil
	case TokenKeyword:
		if tok.Value == "true" {
			return objects.Boolean(true), nil
		}
		if tok.Value == "false" {
			return objects.Boolean(false), nil
		}
		if tok.Value == "null" {
			return nil, nil
		}
		return nil, fmt.Errorf("unexpected keyword: %s", tok.Value)
	}

	return nil, fmt.Errorf("unexpected token: %v (%s)", tok.Type, tok.Value)
}

func (p *Parser) parseArray() (objects.Array, error) {
	var arr objects.Array
	for {
		peek, err := p.peekTokenN(1)
		if err != nil {
			return nil, err
		}
		if peek.Type == TokenCloseBracket {
			p.nextToken()
			break
		}
		obj, err := p.ParseObject()
		if err != nil {
			return nil, err
		}
		arr = append(arr, obj)
	}
	return arr, nil
}

func (p *Parser) parseDictionary() (objects.Dictionary, error) {
	dict := make(objects.Dictionary)
	for {
		peek, err := p.peekTokenN(1)
		if err != nil {
			return nil, err
		}
		if peek.Type == TokenCloseDict {
			p.nextToken()
			break
		}
		keyTok, err := p.nextToken()
		if err != nil || keyTok.Type != TokenName {
			return nil, fmt.Errorf("expected name key in dictionary, got %v", keyTok)
		}
		val, err := p.ParseObject()
		if err != nil {
			return nil, err
		}
		dict[keyTok.Value] = val
	}
	return dict, nil
}

func (p *Parser) nextToken() (Token, error) {
	if len(p.peeks) > 0 {
		tok := p.peeks[0]
		p.peeks = p.peeks[1:]
		return tok, nil
	}
	return p.lexer.NextToken()
}

func (p *Parser) peekTokenN(n int) (*Token, error) {
	for len(p.peeks) < n {
		tok, err := p.lexer.NextToken()
		if err != nil {
			return nil, err
		}
		p.peeks = append(p.peeks, tok)
	}
	return &p.peeks[n-1], nil
}

func (p *Parser) skipEndObj() {
	p.skipKeyword("endobj")
}

func (p *Parser) skipKeyword(keyword string) {
	peek, err := p.peekTokenN(1)
	if err == nil && peek.Type == TokenKeyword && peek.Value == keyword {
		p.nextToken()
	}
}
