package gosql

func parseCreateTableStatement(tokens []*token, initialCursor uint, delimiter token) (*CreateTableStatement, uint, bool) {
	cursor := initialCursor
	if !expectToken(tokens, cursor, tokenFromKeyword(createKeyword)) {
		helpMessage(tokens, cursor, "expected CREATE")
		return nil, initialCursor, false
	}
	cursor++

	if !expectToken(tokens, cursor, tokenFromKeyword(tableKeyword)) {
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

func parseColumnDefinitions(tokens []*token, initialCursor uint, delimiter token) (*[]*columnDefinition, uint, bool) {
	cursor := initialCursor

	cds := []*columnDefinition{}
	for {
		if cursor >= uint(len(tokens)) {
			return nil, initialCursor, false
		}

		current := tokens[cursor]
		if delimiter.equals(current) {
			break
		}

		if len(cds) > 0 {
			if !expectToken(tokens, cursor, tokenFromSymbol(commaSymbol)) {
				helpMessage(tokens, cursor, "Expected ,")
				return nil, initialCursor, false
			}

			cursor++
		}

		name, newCursor, ok := parseToken(tokens, cursor, identifierKind)
		if !ok {
			helpMessage(tokens, cursor, "expected column name")
			return nil, initialCursor, false
		}
		cursor = newCursor

		dataType, newCursor, ok := parseToken(tokens, cursor, keywordKind)
		if !ok {
			helpMessage(tokens, cursor, "expected column type")
			return nil, initialCursor, false
		}
		cursor = newCursor

		cds = append(cds, &columnDefinition{name: *name, datatype: *dataType})
	}

	return &cds, cursor, true
}
