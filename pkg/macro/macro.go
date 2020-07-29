package macro

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var macroExpr = regexp.MustCompile(`\$(\w+)|\[\[([\s\S]+?)(?::(\w+))?\]\]|\${(\w+)(?:\.([^:^\}]+))?(?::(\w+))?}/g`)

var errorProcessorNotFound = errors.New("processor not found")
var processors = map[string]Processor{}

func Register(name string, processor Processor) {
	processors[name] = processor
}

func Interpolate(text string, scope *Scope) (string, error) {
	text = normalizeString(text)
	macrosList := macroExpr.FindAllString(text, 100)
	globalCursor := 0

	for _, macrosDef := range macrosList {
		defStart := strings.Index(text, macrosDef)
		defEnd := defStart + len(macrosDef)

		if strings.HasPrefix(macrosDef, "$__") && text[defEnd] == '(' {
			// Macros
			processorCallable, hasProcessor := processors[macrosDef[3:]]
			if hasProcessor == false {
				return text, errorProcessorNotFound
			}

			globalCursor = defEnd + 1 // skip '('
			args, err := readMacroArgs(text, &globalCursor)
			if err != nil {
				return text, err
			}

			processed, err := processorCallable(args, scope)
			if err != nil {
				return text, err
			}

			text = fmt.Sprintf("%s%s%s", text[0:defStart], processed , text[globalCursor:])
		}
	}
	return text, nil
}

func readMacroArgs(str string, readFrom *int) ([]string, error) {
	args := make([]string, 0)
	strLen := len(str)
	isString := false
	foundClosePar := false

	if *readFrom > strLen {
		return args, errors.New("start index out of string length")
	}

	for i := *readFrom ; i < strLen ; i++ {
		if str[i] == '\'' {
			isString = !isString
		}

		// Args separator found
		if str[i] == ',' && isString == false  {
			args = append(args, strings.TrimSpace(str[*readFrom:i]))
			*readFrom = i + 1 // skip `,`
			continue
		}

		// End of line
		if str[i] == ')' && isString == false  {
			args = append(args, strings.TrimSpace(str[*readFrom:i]))
			foundClosePar = true
			*readFrom = i + 1 // skip ')'
			break
		}
	}

	if foundClosePar == false {
		return args, errors.New("not found closed bracket ')'")
	}
	return args, nil
}

// Removes \r\n from string
func normalizeString(s string) string {
	re := regexp.MustCompile(`\r?\n`)
	return re.ReplaceAllString(s, " ")
}

