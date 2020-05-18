package gosql

import (
	"strings"
)

func lexKeyword(source string, ic cursor) (*token, cursor, bool) {
	cur := ic

	keywords := []keyword{
		selectKeyword,
		insertKeyword,
		valuesKeyword,
		tableKeyword,
		createKeyword,
		dropKeyword,
		whereKeyword,
		fromKeyword,
		intoKeyword,
		textKeyword,
		boolKeyword,
		intKeyword,
		andKeyword,
		orKeyword,
		asKeyword,
		trueKeyword,
		falseKeyword,
		uniqueKeyword,
		indexKeyword,
		onKeyword,
		primarykeyKeyword,
		nullKeyword,
	}

	var options []string
	for _, k := range keywords {
		options = append(options, string(k))
	}

	match := longestMatch(source, ic, options)
	if match == "" {
		return nil, ic, false
	}

	cur.pointer = ic.pointer + uint(len(match))
	cur.loc.col = ic.loc.col + uint(len(match))

	kind := keywordKind
	if match == string(trueKeyword) || match == string(falseKeyword) {
		kind = boolKind
	}

	if match == string(nullKeyword) {
		kind = nullKind
	}

	return &token{
		value: match,
		kind:  kind,
		loc:   ic.loc,
	}, cur, true
}

// longestMatch iterates through a source string starting at the given
// cursor to find the longest matching substring among the provided
// options
func longestMatch(source string, ic cursor, options []string) string {
	var value []byte
	var skipList []int
	var match string

	cur := ic
	for cur.pointer < uint(len(source)) {
		value = append(value, strings.ToLower(string(source[cur.pointer]))...)
		cur.pointer++

	match:
		for i, option := range options {
			for _, skip := range skipList {
				if i == skip {
					continue match
				}
			}

			// Deal with cases like INT vs INTO
			if option == string(value) {
				skipList = append(skipList, i)
				if len(option) > len(match) {
					match = option
				}
				continue
			}

			sharesPrefix := string(value) == option[:cur.pointer-ic.pointer]
			tooLong := len(value) > len(option)
			if tooLong || !sharesPrefix {
				skipList = append(skipList, i)
			}
		}

		if len(skipList) == len(options) {
			break
		}
	}

	return match
}

type keyword string

const (
	selectKeyword     keyword = "select"
	fromKeyword       keyword = "from"
	asKeyword         keyword = "as"
	tableKeyword      keyword = "table"
	createKeyword     keyword = "create"
	dropKeyword       keyword = "drop"
	insertKeyword     keyword = "insert"
	intoKeyword       keyword = "into"
	valuesKeyword     keyword = "values"
	intKeyword        keyword = "int"
	textKeyword       keyword = "text"
	boolKeyword       keyword = "boolean"
	whereKeyword      keyword = "where"
	andKeyword        keyword = "and"
	orKeyword         keyword = "or"
	trueKeyword       keyword = "true"
	falseKeyword      keyword = "false"
	uniqueKeyword     keyword = "unique"
	indexKeyword      keyword = "index"
	onKeyword         keyword = "on"
	primarykeyKeyword keyword = "primary key"
	nullKeyword       keyword = "null"
)
