package gosql

func lexSymbol(source string, ic cursor) (*token, cursor, bool) {
	c := source[ic.pointer]
	cur := ic
	cur.loc.col++
	cur.pointer++

	switch c {
	// Syntax that should be thrown away
	case '\n':
		cur.loc.line++
		cur.loc.col = 0
		fallthrough
	case '\t':
		fallthrough
	case ' ':
		return nil, cur, true
	// Syntax that should be kept
	case ',':
		fallthrough
	case '(':
		fallthrough
	case ')':
		fallthrough
	case ';':
		fallthrough
	case '*':
		break
	// Unknown character
	default:
		return nil, ic, false
	}

	return &token{
		value: string(c),
		loc:   ic.loc,
		kind:  symbolKind,
	}, cur, true
}
