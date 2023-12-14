package main

import (
	"fmt"
	"os"
	"reflect"
	"unicode"
)

type JsonValue struct {
	Identifier string
	Value      any
	Type       reflect.Type
}

type TokenType int

const (
	TokenOpenBrace TokenType = iota
	TokenCloseBrace
	TokenComma
	TokenColon
	TokenOpenBracket
	TokenCloseBracket
	TokenNumber
	TokenString
)

var TokenMap = map[rune]TokenType{
	'{': TokenOpenBrace,
	'}': TokenCloseBrace,
	',': TokenComma,
	':': TokenColon,
	'[': TokenOpenBracket,
	']': TokenCloseBracket,
}

type JsonToken struct {
	Type  TokenType
	Value string
}

func getJsonToken(reader *Reader) (JsonToken, error) {

	for ; reader.At < reader.Len && isJsonWhiteSpace(reader.Buffer[reader.At]); reader.At++ {
	}
	var token JsonToken
	if reader.At >= reader.Len {
		return token, fmt.Errorf("end")
	}
	r := reader.Buffer[reader.At]
	reader.At++
	// Small one char token
	if tokenType, exist := TokenMap[r]; exist {
		token.Type = tokenType
		token.Value = string(r)
	} else if r == '"' {
		// find the closing quote
		j := reader.At
		for ; j < reader.Len && reader.Buffer[j] != '"'; j++ {
		}

		if reader.Buffer[j] != '"' {
			return token, fmt.Errorf("missing '\"'")
		}
		token.Type = TokenString
		if j-reader.At == 0 {
			token.Value = ""
		} else {
			token.Value = string(reader.Buffer[reader.At:j])
		}
		reader.At = j + 1
	} else if (r >= '0' && r <= '9') || r == '-' {
		// find the end of the number
		j := reader.At
		for ; j < reader.Len && isInNumber(reader.Buffer[j]); j++ {
		}

		token.Type = TokenNumber
		token.Value = string(reader.Buffer[reader.At-1 : j])
		reader.At = j
	} else {
		return token, fmt.Errorf("not supported")
	}

	return token, nil
}

func isJsonWhiteSpace(r rune) bool {
	return unicode.IsSpace(r) || unicode.IsControl(r)
}

func isInNumber(r rune) bool {
	return (r >= '0' && r <= '9') || r == '.'
}

type Reader struct {
	Buffer []rune
	At     int
	Len    int
}

func main() {
	file, err := os.ReadFile("data.json")

	p(err)

	buffer := []rune(string(file))

	reader := Reader{
		Buffer: buffer,
		At:     0,
		Len:    len(buffer),
	}

	for {
		jsonToken, err := getJsonToken(&reader)
		p(err)

		fmt.Println(jsonToken)
	}

}

func p(err error) {
	if err != nil {
		panic(err)
	}
}
