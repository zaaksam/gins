package orm

import (
	"bytes"
	"unicode/utf8"
)

// 参考 easyjson/jwriter/writer.go 实现
const chars = "0123456789abcdef"

// 参考 easyjson/jwriter/writer.go 实现
func getTable(falseValues ...int) [128]bool {
	table := [128]bool{}

	for i := 0; i < 128; i++ {
		table[i] = true
	}

	for _, v := range falseValues {
		table[v] = false
	}

	return table
}

// 参考 easyjson/jwriter/writer.go 实现
var (
	htmlEscapeTable   = getTable(0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, '"', '&', '<', '>', '\\')
	htmlNoEscapeTable = getTable(0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, '"', '\\')
)

// 参考 easyjson/jwriter/writer.go 实现
func stringToJSONValue(str string, noEscapeHTMLOpt ...bool) (body []byte) {
	var jsonValue bytes.Buffer

	escapeTable := &htmlEscapeTable
	if len(noEscapeHTMLOpt) == 1 && noEscapeHTMLOpt[0] {
		escapeTable = &htmlNoEscapeTable
	}

	jsonValue.WriteByte('"')

	p := 0
	for i := 0; i < len(str); {
		c := str[i]

		if c < utf8.RuneSelf {
			if escapeTable[c] {
				i++
				continue
			}

			jsonValue.WriteString(str[p:i])
			switch c {
			case '\t':
				jsonValue.WriteString(`\t`)
			case '\r':
				jsonValue.WriteString(`\r`)
			case '\n':
				jsonValue.WriteString(`\n`)
			case '\\':
				jsonValue.WriteString(`\\`)
			case '"':
				jsonValue.WriteString(`\"`)
			default:
				jsonValue.WriteString(`\u00`)
				jsonValue.WriteByte(chars[c>>4])
				jsonValue.WriteByte(chars[c&0xf])
			}

			i++
			p = i
			continue
		}

		runeValue, runeWidth := utf8.DecodeRuneInString(str[i:])
		if runeValue == utf8.RuneError && runeWidth == 1 {
			jsonValue.WriteString(str[p:i])
			jsonValue.WriteString(`\ufffd`)
			i++
			p = i
			continue
		}

		if runeValue == '\u2028' || runeValue == '\u2029' {
			jsonValue.WriteString(str[p:i])
			jsonValue.WriteString(`\u202`)
			jsonValue.WriteByte(chars[runeValue&0xf])
			i += runeWidth
			p = i
			continue
		}
		i += runeWidth
	}
	jsonValue.WriteString(str[p:])

	jsonValue.WriteByte('"')

	return jsonValue.Bytes()
}
