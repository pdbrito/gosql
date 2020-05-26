package gosql

import (
	"errors"
	"fmt"
)

func tokenFromKeyword(k keyword) token {
	return token{kind: keywordKind, value: string(k)}
}

func tokenFromSymbol(s symbol) token {
	return token{kind: symbolKind, value: string(s)}
}

func expectToken(tokens []*token, cursor uint, t token) bool {
	if cursor >= uint(len(tokens)) {
		return false
	}

	return t.equals(tokens[cursor])
}

func helpMessage(tokens []*token, cursor uint, msg string) {
	var c *token
	if cursor < uint(len(tokens)) {
		c = tokens[cursor]
	} else {
		c = tokens[cursor-1]
	}

	fmt.Printf("[%d,%d]: %s, got: %s\n", c.loc.line, c.loc.col, msg, c.value)
}

func Parse(source string) (*Ast, error) {
	tokens, err := lex(source)
	if err != nil {
		return nil, err
	}
	a := Ast{}
	cursor := uint(0)
	for cursor < uint(len(tokens)) {
		stmt, newCursor, ok := parseStatement(tokens, cursor, tokenFromSymbol(semicolonSymbol))
		if !ok {
			helpMessage(tokens, cursor, "Expected statement")
			return nil, errors.New("failed to parse, expected statement")
		}
		cursor = newCursor

		a.Statements = append(a.Statements, stmt)

		atLeastOneSemicolon := false
		for expectToken(tokens, cursor, tokenFromSymbol(semicolonSymbol)) {
			cursor++
			atLeastOneSemicolon = true
		}

		if !atLeastOneSemicolon {
			helpMessage(tokens, cursor, "Expected semicolon delimiter between statements")
			return nil, errors.New("missing semicolon between statements")
		}
	}

	return &a, nil
}

func parseStatement(tokens []*token, initialCursor uint, delimiter token) (*Statement, uint, bool) {
	cursor := initialCursor

	// Look for SELECT statement
	semicolonToken := tokenFromSymbol(semicolonSymbol)
	slct, newCursor, ok := parseSelectStatement(tokens, cursor, semicolonToken)
	if ok {
		return &Statement{
			Kind:            SelectKind,
			SelectStatement: slct,
		}, newCursor, true
	}

	// Look for INSERT statement
	inst, newCursor, ok := parseInsertStatement(tokens, cursor, semicolonToken)
	if ok {
		return &Statement{
			Kind:            InsertKind,
			SelectStatement: slct,
		}, newCursor, true
	}

	// Look for CREATE statment
	crtTbl, newCursor, ok := parseCreateTableStatement(tokens, cursor, semicolonToken)
	if ok {
		return &Statement{
			Kind:                 CreateTableKind,
			CreateTableStatement: crtTbl,
		}, newCursor, true
	}

	return nil, initialCursor, false
}

func parseToken(tokens []*token, initialCursor uint, kind tokenKind) (*token, uint, bool) {
	cursor := initialCursor

	if cursor >= uint(len(tokens)) {
		return nil, initialCursor, false
	}

	current := tokens[cursor]
	if current.kind == kind {
		return current, cursor + 1, true
	}

	return nil, initialCursor, false
}

func parseExpressions(tokens []*token, initialCursor uint, delimiters []token) (*[]*expression, uint, bool) {
	cursor := initialCursor

	exps := []*expression{}

outer:
	for {
		if cursor >= uint(len(tokens)) {
			return nil, initialCursor, false
		}

		current := tokens[cursor]
		for _, delimiter := range delimiters {
			if delimiter.equals(current) {
				break outer
			}
		}

		if len(exps) > 0 {
			if !expectToken(tokens, cursor, tokenFromSymbol(commaSymbol)) {
				helpMessage(tokens, cursor, "Expected comma")
				return nil, initialCursor, false
			}

			cursor++
		}

		exp, newCursor, ok := parseExpression(tokens, cursor, tokenFromSymbol(commaSymbol))
		if !ok {
			helpMessage(tokens, cursor, "Expected expression")
			return nil, initialCursor, false
		}

		cursor = newCursor
		exps = append(exps, exp)
	}

	return &exps, cursor, true
}

func parseExpression(tokens []*token, initialCursor uint, _ token) (*expression, uint, bool) {
	cursor := initialCursor

	kinds := []tokenKind{identifierKind, numericKind, stringKind}
	for _, kind := range kinds {
		t, newCursor, ok := parseToken(tokens, cursor, kind)
		if ok {
			return &expression{
				literal: t,
				kind:    literalKind,
			}, newCursor, true
		}
	}

	return nil, initialCursor, false
}
