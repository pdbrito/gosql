package gosql

func lexSymbol(source string, ic cursor) (*token, cursor, bool) {
	c := source[ic.pointer]
	cur := ic
	cur.pointer++
	cur.loc.col++

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
	}

	symbols := []symbol{
		eqSymbol,
		neqSymbol,
		neqSymbol2,
		ltSymbol,
		lteSymbol,
		gtSymbol,
		gteSymbol,
		concatSymbol,
		plusSymbol,
		commaSymbol,
		leftParenSymbol,
		rightParenSymbol,
		semicolonSymbol,
		asteriskSymbol,
	}

	var options []string
	for _, s := range symbols {
		options = append(options, string(s))
	}

	match := longestMatch(source, ic, options)

	if match == "" {
		return nil, ic, false
	}

	cur.pointer = ic.pointer + uint(len(match))
	cur.loc.col = ic.loc.col + uint(len(match))

	return &token{
		value: match,
		loc:   ic.loc,
		kind:  symbolKind,
	}, cur, true
}

type symbol string

const (
	semicolonSymbol  symbol = ";"
	asteriskSymbol   symbol = "*"
	commaSymbol      symbol = ","
	leftParenSymbol  symbol = "("
	rightParenSymbol symbol = ")"
	eqSymbol         symbol = "="
	neqSymbol        symbol = "<>"
	neqSymbol2       symbol = "!="
	concatSymbol     symbol = "||"
	plusSymbol       symbol = "+"
	ltSymbol         symbol = "<"
	lteSymbol        symbol = "<="
	gtSymbol         symbol = ">"
	gteSymbol        symbol = ">="
)
