package gosql

func parseCreateTableStatement(tokens []*token, initialCursor uint, delimiter token) (*CreateTableStatement, uint, bool) {
	cursor := initialCursor
	if !expectToken(tokens, cursor, tokenFromKeyword(createKeyword)) {
		helpMessage(tokens, cursor, "expected CREATE")
		return nil, initialCursor, false
	}
	cursor++

	name, newCursor, ok := parseToken(tokens, cursor, identifierKind)
	if !ok {
		helpMessage(tokens, cursor, "expected table name")
		return nil, initialCursor, false
	}
	cursor = newCursor

	if !expectToken(tokens, cursor, tokenFromSymbol(leftParenSymbol)) {
		helpMessage(tokens, cursor, "expected (")
		return nil, initialCursor, false
	}
	cursor++

	cols, newCursor, ok := parseColumnDefinitions(tokens, cursor, tokenFromSymbol(rightParenSymbol))
	if !ok {
		return nil, initialCursor, false
	}
	cursor = newCursor

	if !expectToken(tokens, cursor, tokenFromSymbol(rightParenSymbol)) {
		helpMessage(tokens, cursor, "expected )")
		return nil, initialCursor, false
	}
	cursor++

	return &CreateTableStatement{
		name: *name,
		cols: cols,
	}, cursor, true
}
