package wirtokenizer

import (
	"fmt"
	"strings"

	"github.com/phillip-england/wir/internal/runelexer"
	"github.com/phillip-england/wir/internal/wherr"
)

type TokenType string

const (
	TokenTypeRawText                 = "RAW_TEXT"
	TokenTypeString                  = "STRING"
	TokenTypeTagInfo                 = "TAG_INFO"
	TokenTypeAtDirective             = "AT_DIRECTIVE"
	TokenTypeCurlyBraceOpen          = "CURLY_BRACE_OPEN"
	TokenTypeCurlyBraceClose         = "CURLY_BRACE_CLOSE"
	TokenTypeDollarSignInterpolation = "DOLLAR_SIGN_INTERPOLATION"

	TokenTypeHTMLTagInfoStart    = "HTML_TAG_INFO_START"
	TokenTypeHTMLTagInfoEnd      = "HTML_TAG_INFO_END"
	TokenTypeHTMLTagName         = "HTML_TAG_NAME"
	TokenTypeHTMLCurlyBraceOpen  = "HTML_CURLY_BRACE_OPEN"
	TokenTypeHTMLCurlyBraceClose = "HTML_CURLY_BRACE_CLOSE"

	TokenTypeHTMLAttrKey          = "HTML_ATTR_KEY"
	TokenTypeHTMLAttrEqualSign    = "HTML_ATTR_EQUAL_SIGN"
	TokenTypeHTMLAttrValue        = "HTML_ATTR_VALUE"
	TokenTypeHTMLAttrValuePartial = "HTML_ATTR_VALUE_PARTIAL"

	TokenTypeStringStart   = "STRING_START"
	TokenTypeStringEnd     = "STRING_END"
	TokenTypeStringContent = "STRING_CONTENT"

	TokenTypeDollarSignInterpolationOpen      = "DOLLAR_SIGN_INTERPOLATION_OPEN"
	TokenTypeDollarSignInterpolationClose     = "DOLLAR_SIGN_INTERPOLATION_CLOSE"
	TokenTypeDollarSignInterpolationValue     = "DOLLAR_SIGN_INTERPOLATION_VALUE"
	TokenTypeDollarSignInterpolationSemiColon = "DOLLAR_SIGN_INTERPOLATION_SEMICOLON"
	TokenTypeDollarSignInterpolationType      = "DOLLAR_SIGN_INTERPOLATION_TYPE"

	TokenTypeAtDirectiveStart            = "AT_DIRECTIVE_START"
	TokenTypeAtDirectiveName             = "AT_DIRECTIVE_NAME"
	TokenTypeAtDirectiveParenthesisOpen  = "AT_DIRECTIVE_PARENTHESIS_OPEN"
	TokenTypeAtDirectiveParenthesisClose = "AT_DIRECTIVE_PARENTHESIS_CLOSE"
	TokenTypeAtDirectiveParamValue       = "AT_DIRECTIVE_PARAM_VALUE"
	TokenTypeAtDirectiveSemiColon        = "AT_DIRECTIVE_SEMICOLON"
	TokenTypeAtDirectiveParamType        = "AT_DIRECTIVE_PARAM_TYPE"


	TokenTypeEndOfFile = "END_OF_FILE"
)

type Token struct {
	t    TokenType
	text string
}

func (t Token) Str() string {
	return fmt.Sprintf("%s:%s", t.t, t.text)
}

func TokenManyFromStr(s string) ([]Token, error) {
	var toks []Token
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		l := runelexer.NewRuneLexer[Token](line)
		foundColon := false
		l.Iter(func(ch string, pos int) bool {
			if ch == ":" {
				foundColon = true
				return false
			}
			return true
		})
		if !foundColon {
			return toks, wherr.Err(wherr.Here(), "token string is malformed on line %d: %s", i+1, line)
		}
		t := l.PullFromStart()
		l.Next()
		text := l.PullFromEnd()
		tok := Token{
			t: TokenType(t),
			text: text,
		}
		toks = append(toks, tok)
	}
	return toks, nil
} 