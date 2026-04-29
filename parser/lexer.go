package parser

import (
	"fmt"
	"strings"
)

const spacesPerIndent = 4

// Takes a raw Waymark source string and converts it to a flat slice of tokens
func Tokenize(src string) ([]Token, error) {
	var tokens []Token
	lines := strings.Split(src, "\n")

	for i, raw := range lines {
		// Strip \r for Windows line endings
		raw = strings.TrimRight(raw, "\r")

		trimmed := strings.TrimSpace(raw)

		// Skip blanks and comments
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		indent, err := countIndent(raw)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", i+1, err)
		}

		tok := Token{Indent: indent, Line: i + 1}

		switch {
		case strings.HasPrefix(trimmed, "scene "):
			tok.Type, tok.Value = TOK_SCENE, after(trimmed, "scene ")
		case strings.HasPrefix(trimmed, "var "):
			tok.Type, tok.Value = TOK_VAR, after(trimmed, "var ")
		case strings.HasPrefix(trimmed, "local "):
			tok.Type, tok.Value = TOK_LOCAL, after(trimmed, "local ")
		case strings.HasPrefix(trimmed, "keep "):
			tok.Type, tok.Value = TOK_KEEP, after(trimmed, "keep ")
		case strings.HasPrefix(trimmed, "set "):
			tok.Type, tok.Value = TOK_SET, after(trimmed, "set ")
		case strings.HasPrefix(trimmed, "input "):
			tok.Type, tok.Value = TOK_INPUT, after(trimmed, "input ")
		case trimmed == "choice":
			tok.Type = TOK_CHOICE
		case trimmed == "go" || strings.HasPrefix(trimmed, "go "):
			tok.Type, tok.Value = TOK_GO, after(trimmed, "go ")
		case trimmed == "call" || strings.HasPrefix(trimmed, "call "):
			tok.Type, tok.Value = TOK_CALL, after(trimmed, "call ")
		case trimmed == "finish":
			tok.Type = TOK_FINISH
		case trimmed == "end":
			tok.Type = TOK_END
		case strings.HasPrefix(trimmed, `"`):
			// NOTE: Quoted strings in a `choice` block are emitted as DIALOGUE, not BRANCH
			closeQuote := strings.Index(trimmed[1:], `"`) + 1
			label := trimmed[1:closeQuote]
			rest := strings.TrimSpace(trimmed[closeQuote+1:])

			if strings.HasPrefix(rest, "when ") {
				tok.Type = TOK_BRANCH
				tok.Value = label
				tok.Condition = after(rest, "when ")
			} else {
				tok.Type = TOK_DIALOGUE
				tok.Value = label
			}
		default:
			tok.Type, tok.Value = TOK_TEXT, trimmed
		}

		tokens = append(tokens, tok)
	}

	return tokens, nil
}

// Count the number of indents in a line. Accepts tabs and spaces, for now
// defined static in spacesPerIndent
func countIndent(line string) (int, error) {
	tabs, spaces := 0, 0

	for _, ch := range line {
		switch ch {
		case '\t':
			if spaces > 0 {
				return 0, fmt.Errorf("mixed tabs and spaces in indent")
			}
			tabs++
		case ' ':
			if tabs > 0 {
				return 0, fmt.Errorf("mixed tabs and spaces in indent")
			}
			spaces++
		default:
			goto done
		}
	}

done:
	if spaces > 0 {
		if spaces%spacesPerIndent != 0 {
			return 0, fmt.Errorf(
				"indent of %d spaces is not a multiple of %d",
				spaces, spacesPerIndent,
			)
		}
		return spaces / spacesPerIndent, nil
	}
	return tabs, nil
}

// Wrapper for trimming prefix from line
func after(s, prefix string) string {
	return strings.TrimSpace(strings.TrimPrefix(s, prefix))
}
