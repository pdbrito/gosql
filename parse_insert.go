package gosql

func parseInsertStatement(tokens []*token, initialCursor uint, delimiter token) (*InsertStatement, uint, bool) {
	cursor := initialCursor
	if !expectToken(tokens, cursor, tokenFromKeyword(insertKeyword)) {
		helpMessage(tokens, cursor, "expected INSERT")
		return nil, initialCursor, false
	}
	cursor++

	if !expectToken(tokens, cursor, tokenFromKeyword(intoKeyword)) {
		helpMessage(tokens, cursor, "expected INTO")
		return nil, initialCursor, false
	}
	cursor++

	table, newCursor, ok := parseToken(tokens, cursor, identifierKind)
	if !ok {
		helpMessage(tokens, cursor, "expected table name")
		return nil, initialCursor, false
	}
	cursor = newCursor

	if !expectToken(tokens, cursor, tokenFromKeyword(valuesKeyword)) {
		helpMessage(tokens, cursor, "expected VALUES")
		return nil, initialCursor, false
	}
	cursor++

	if !expectToken(tokens, cursor, tokenFromSymbol(leftParenSymbol)) {
		helpMessage(tokens, cursor, "expected (")
		return nil, initialCursor, false
	}
	cursor++

	values, newCursor, ok := parseExpressions(tokens, cursor, tokenFromSymbol(rightParenSymbol))
	if !ok {
		helpMessage(tokens, cursor, "expected one or more comma separated values")
		return nil, initialCursor, false
	}
	cursor = newCursor

	if !expectToken(tokens, cursor, tokenFromSymbol(rightParenSymbol)) {
		helpMessage(tokens, cursor, "expected )")
		return nil, initialCursor, false
	}
	cursor++

	return &InsertStatement{
		table:  *table,
		values: values,
	}, cursor, true
}
