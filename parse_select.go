package gosql

// SELECT [ident [, ...]] [FROM ident]
func parseSelectStatement(tokens []*token, initialCursor uint, delimiter token) (*SelectStatement, uint, bool) {
	cursor := initialCursor
	if !expectToken(tokens, cursor, tokenFromKeyword(selectKeyword)) {
		return nil, initialCursor, false
	}
	cursor++

	slct := SelectStatement{}

	item, newCursor, ok := parseSelectItem(tokens, cursor, []token{tokenFromKeyword(fromKeyword), delimiter})
	if !ok {
		return nil, initialCursor, false
	}

	slct.item = item
	cursor = newCursor

	if expectToken(tokens, cursor, tokenFromKeyword(fromKeyword)) {
		cursor++

		from, newCursor, ok := parseFromItem(tokens, cursor, delimiter)
		if !ok {
			helpMessage(tokens, cursor, "Expected FROM item")
			return nil, initialCursor, false
		}

		slct.from = from
		cursor = newCursor
	}

	return &slct, cursor, true
}

// expression [AS ident] [, ...]
func parseSelectItem(tokens []*token, initialCursor uint, delimiters []token) (*[]*selectItem, uint, bool) {
	cursor := initialCursor

	s := []*selectItem{}

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

		if len(s) > 0 {
			if !expectToken(tokens, cursor, tokenFromSymbol(commaSymbol)) {
				helpMessage(tokens, cursor, "Expected ,")
				return nil, initialCursor, false
			}

			cursor++
		}

		var si selectItem
		if expectToken(tokens, cursor, tokenFromSymbol(asteriskSymbol)) {
			si = selectItem{asterisk: true}
			cursor++
			s = append(s, &si)
			continue
		}

		exp, newCursor, ok := parseExpression(tokens, cursor, tokenFromSymbol(commaSymbol))
		if !ok {
			helpMessage(tokens, cursor, "Expected expression")
			return nil, initialCursor, false
		}

		cursor = newCursor
		si.exp = exp

		if expectToken(tokens, cursor, tokenFromKeyword(asKeyword)) {
			cursor++

			id, newCursor, ok := parseToken(tokens, cursor, identifierKind)
			if !ok {
				helpMessage(tokens, cursor, "Expected identifier after AS")
				return nil, initialCursor, false
			}

			cursor = newCursor
			si.as = id
		}
		s = append(s, &si)
	}
	return &s, cursor, true
}

func parseFromItem(tokens []*token, initialCursor uint, _ token) (*fromItem, uint, bool) {
	ident, newCursor, ok := parseToken(tokens, initialCursor, identifierKind)
	if !ok {
		return nil, initialCursor, false
	}

	return &fromItem{table: ident}, newCursor, true
}
