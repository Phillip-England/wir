package wirtokenizer

import (
	"strings"

	"github.com/phillip-england/wir/internal/runelexer"
)

type WirTokenType string

const (
	WirTokenTypeRawText                 = "RAW_TEXT"
	WirTokenTypeString                  = "STRING"
	WirTokenTypeTagInfo                 = "TAG_INFO"
	WirTokenTypeAtDirective             = "AT_DIRECTIVE"
	WirTokenTypeCurlyBraceOpen          = "CURLY_BRACE_OPEN"
	WirTokenTypeCurlyBraceClose         = "CURLY_BRACE_CLOSE"
	WirTokenTypeDollarSignInterpolation = "DOLLAR_SIGN_INTERPOLATION"

	WirTokenTypeHTMLTagInfoStart    = "HTML_TAG_INFO_START"
	WirTokenTypeHTMLTagInfoEnd      = "HTML_TAG_INFO_END"
	WirTokenTypeHTMLTagName         = "HTML_TAG_NAME"
	WirTokenTypeHTMLCurlyBraceOpen  = "HTML_CURLY_BRACE_OPEN"
	WirTokenTypeHTMLCurlyBraceClose = "HTML_CURLY_BRACE_CLOSE"

	WirTokenTypeHTMLAttrKey          = "HTML_ATTR_KEY"
	WirTokenTypeHTMLAttrEqualSign    = "HTML_ATTR_EQUAL_SIGN"
	WirTokenTypeHTMLAttrValue        = "HTML_ATTR_VALUE"
	WirTokenTypeHTMLAttrValuePartial = "HTML_ATTR_VALUE_PARTIAL"

	WirTokenTypeStringStart   = "STRING_START"
	WirTokenTypeStringEnd     = "STRING_END"
	WirTokenTypeStringContent = "STRING_CONTENT"

	WirTokenTypeDollarSignInterpolationOpen      = "DOLLAR_SIGN_INTERPOLATION_OPEN"
	WirTokenTypeDollarSignInterpolationClose     = "DOLLAR_SIGN_INTERPOLATION_CLOSE"
	WirTokenTypeDollarSignInterpolationValue     = "DOLLAR_SIGN_INTERPOLATION_VALUE"
	WirTokenTypeDollarSignInterpolationSemiColon = "DOLLAR_SIGN_INTERPOLATION_SEMICOLON"
	WirTokenTypeDollarSignInterpolationType      = "DOLLAR_SIGN_INTERPOLATION_TYPE"

	WirTokenTypeAtDirectiveStart            = "AT_DIRECTIVE_START"
	WirTokenTypeAtDirectiveName             = "AT_DIRECTIVE_NAME"
	WirTokenTypeAtDirectiveParenthesisOpen  = "AT_DIRECTIVE_PARENTHESIS_OPEN"
	WirTokenTypeAtDirectiveParenthesisClose = "AT_DIRECTIVE_PARENTHESIS_CLOSE"
	WirTokenTypeAtDirectiveParamValue       = "AT_DIRECTIVE_PARAM_VALUE"
	WirTokenTypeAtDirectiveSemiColon        = "AT_DIRECTIVE_SEMICOLON"
	WirTokenTypeAtDirectiveParamType        = "AT_DIRECTIVE_PARAM_TYPE"

	WirTokenTypeEndOfFile = "END_OF_FILE"
)

type WirToken struct {
	t    WirTokenType
	text string
}

type WirTokenizer struct {
	lexer *runelexer.RuneLexer[WirToken]
}

func NewFromString(s string) (*WirTokenizer, error) {
	s = strings.TrimSpace(s)
	l := runelexer.NewRuneLexer[WirToken](s)
	err := tokenizeWir(l)
	if err != nil {
		return &WirTokenizer{}, err
	}
	return &WirTokenizer{
		lexer: l,
	}, nil
}

func (t *WirTokenizer) TokStr() string {
	s := ""
	if t.lexer.TokenLen() == 0 {
		return ""
	}
	t.lexer.TokenIter(func(token WirToken, index int) bool {
		s += string(token.t) + ":" + token.text + "\n"
		return true
	})
	return s
}

type WirTokenState = int

const (
	WirTokenStateInit = iota
)

func tokenizeWir(l *runelexer.RuneLexer[WirToken]) error {
	err := phase1(l)
	if err != nil {
		return err
	}
	err = phase2(l)
	err = phase3(l)
	return nil
}

func phase3(l *runelexer.RuneLexer[WirToken]) error {
	var toks []WirToken
	l.TokenIter(func(tk WirToken, index int) bool {
		switch tk.t {
		default:
			{
				toks = append(toks, tk)
			}
		case WirTokenTypeCurlyBraceOpen:
			{
				toks = append(toks, WirToken{
					t:    WirTokenTypeHTMLCurlyBraceOpen,
					text: tk.text,
				})
			}
		case WirTokenTypeCurlyBraceClose:
			{
				toks = append(toks, WirToken{
					t:    WirTokenTypeHTMLCurlyBraceClose,
					text: tk.text,
				})
			}
		case WirTokenTypeRawText:
			{
				toks = append(toks, WirToken{
					t:    WirTokenTypeHTMLTagName,
					text: tk.text,
				})
			}
		case WirTokenTypeDollarSignInterpolation:
			{
				toks = append(toks, WirToken{
					t:    WirTokenTypeDollarSignInterpolationOpen,
					text: "${",
				})
				s := tk.text
				s = strings.Replace(s, "${", "", 1)
				s = s[0 : len(s)-1]
				innerParts := strings.Split(s, ":")
				for i, v := range innerParts {
					if i%2 == 0 {
						toks = append(toks, WirToken{
							t:    WirTokenTypeDollarSignInterpolationValue,
							text: strings.TrimSpace(v),
						})
						toks = append(toks, WirToken{
							t:    WirTokenTypeDollarSignInterpolationSemiColon,
							text: ":",
						})
					} else {
						toks = append(toks, WirToken{
							t:    WirTokenTypeDollarSignInterpolationType,
							text: strings.TrimSpace(v),
						})
					}
				}
			}
			toks = append(toks, WirToken{
				t:    WirTokenTypeDollarSignInterpolationClose,
				text: "}",
			})
		}
		return true
	})
	l.TokenOverwrite(toks)
	return nil
}

func phase2(l *runelexer.RuneLexer[WirToken]) error {
	var toks []WirToken
	l.TokenIter(func(tk WirToken, index int) bool {
		switch tk.t {
		default:
			{
				toks = append(toks, tk)
			}
		case WirTokenTypeTagInfo:
			{
				l2 := runelexer.NewRuneLexer[WirToken](tk.text)
				l2.Iter(func(ch string, pos int) bool {
					switch ch {
					default:
						{
							l2.Store()
						}
					case "=":
						{
							attrKey := strings.TrimSpace(l2.StoreFlush())
							toks = append(toks, WirToken{
								t:    WirTokenTypeHTMLAttrKey,
								text: attrKey,
							})
							toks = append(toks, WirToken{
								t:    WirTokenTypeHTMLAttrEqualSign,
								text: "=",
							})
						}
					case "<":
						{
							toks = append(toks, WirToken{
								t:    WirTokenTypeHTMLTagInfoStart,
								text: "<",
							})
						}
					case ">":
						{
							attrKey := strings.TrimSpace(l2.StoreFlush())
							toks = append(toks, WirToken{
								t:    WirTokenTypeHTMLAttrKey,
								text: attrKey,
							})
							toks = append(toks, WirToken{
								t:    WirTokenTypeHTMLTagInfoEnd,
								text: ">",
							})
						}
					case "'":
						{
							l2.Mark()
							l2.Next()
							l2.Iter(func(ch2 string, pos int) bool {
								if ch2 == "'" && l2.Peek(-1) != "\\" {
									htmlAttr := l2.PullFromMark()
									brokeAttr := false
									l3 := runelexer.NewRuneLexer[WirToken](htmlAttr)
									l3.Iter(func(ch3 string, pos int) bool {
										switch ch3 {
										default:
											{
												l3.Store()
											}
										case "$":
											if l3.Peek(1) != "{" {
												l3.Store()
												return true
											}
											brokeAttr = true
											attrBit := l3.StoreFlush()
											toks = append(toks, WirToken{
												t:    WirTokenTypeHTMLAttrValuePartial,
												text: attrBit,
											})
											l3.Mark()
											l3.NextUntil("}")
											dollarSignInterpolation := l3.PullFromMark()
											toks = append(toks, WirToken{
												t:    WirTokenTypeDollarSignInterpolation,
												text: dollarSignInterpolation,
											})
										}
										return true
									})
									if brokeAttr {
										toks = append(toks, WirToken{
											t:    WirTokenTypeHTMLAttrValuePartial,
											text: l3.StoreFlush(),
										})
									} else {
										toks = append(toks, WirToken{
											t:    WirTokenTypeHTMLAttrValue,
											text: l3.PullFromMark(),
										})
									}
									return false
								}
								return true
							})
						}
					case "\"":
						{
							l2.Mark()
							l2.Next()
							l2.Iter(func(ch2 string, pos int) bool {
								if ch2 == "\"" && l2.Peek(-1) != "\\" {
									toks = append(toks, WirToken{
										t:    WirTokenTypeHTMLAttrValue,
										text: l2.PullFromMark(),
									})
									return false
								}
								return true
							})
						}
					}
					return true
				})
			}
		case WirTokenTypeString:
			{
				l2 := runelexer.NewRuneLexer[WirToken](tk.text)
				l2.Iter(func(ch string, pos int) bool {
					switch ch {
					default:
						{
							l2.Store()
						}
					case "$":
						{
							if l2.Peek(1) != "{" {
								l2.Store()
								return true
							}
							toks = append(toks, WirToken{
								t:    WirTokenTypeStringContent,
								text: l2.StoreFlush(),
							})
							l2.Mark()
							l2.NextUntil("}")
							toks = append(toks, WirToken{
								t:    WirTokenTypeDollarSignInterpolation,
								text: l2.PullFromMark(),
							})
						}
					case "'":
						{
							if pos == 0 {
								toks = append(toks, WirToken{
									t:    WirTokenTypeStringStart,
									text: "'",
								})
								return true
							}
							if l2.AtEnd() && ch == "'" {
								toks = append(toks, WirToken{
									t:    WirTokenTypeStringContent,
									text: l2.StoreFlush(),
								})
								toks = append(toks, WirToken{
									t:    WirTokenTypeStringEnd,
									text: "'",
								})
								return true
							}
						}
					case "\"":
						{
							if pos == 0 {
								toks = append(toks, WirToken{
									t:    WirTokenTypeStringStart,
									text: "\"",
								})
								return true
							}
							if l2.AtEnd() && ch == "\"" {
								toks = append(toks, WirToken{
									t:    WirTokenTypeStringContent,
									text: l2.StoreFlush(),
								})
								toks = append(toks, WirToken{
									t:    WirTokenTypeStringEnd,
									text: "\"",
								})
								return true
							}
						}
					}
					return true
				})
			}
		case WirTokenTypeAtDirective:
			{
				l2 := runelexer.NewRuneLexer[WirToken](tk.text)
				l2.Iter(func(ch string, pos int) bool {
					switch ch {
					default:
						{
							l2.Store()
						}
					case ")":
						{
							toks = append(toks, WirToken{
								t:    WirTokenTypeAtDirectiveParenthesisClose,
								text: ")",
							})
						}
					case "(":
						{
							directiveName := l2.StoreFlush()
							toks = append(toks, WirToken{
								t:    WirTokenTypeAtDirectiveName,
								text: directiveName,
							})
							toks = append(toks, WirToken{
								t:    WirTokenTypeAtDirectiveParenthesisOpen,
								text: "(",
							})
							l2.Next()
							l2.Mark()
							l2.GoToEnd()
							l2.Prev()
							l2.PullFromMark()
							directiveInputParams := l2.PullFromMark()
							directiveInputParts := strings.Split(directiveInputParams, ":")
							for i, v := range directiveInputParts {
								v = strings.TrimSpace(v)
								if i%2 == 0 {
									toks = append(toks, WirToken{
										t:    WirTokenTypeAtDirectiveParamValue,
										text: v,
									})
									toks = append(toks, WirToken{
										t:    WirTokenTypeAtDirectiveSemiColon,
										text: ":",
									})
								} else {
									toks = append(toks, WirToken{
										t:    WirTokenTypeAtDirectiveParamType,
										text: v,
									})
								}
							}
						}
					case "@":
						{
							toks = append(toks, WirToken{
								t:    WirTokenTypeAtDirectiveStart,
								text: "@",
							})
						}
					}
					return true
				})
			}
		}
		return true
	})
	l.TokenOverwrite(toks)
	return nil
}

func phase1(l *runelexer.RuneLexer[WirToken]) error {
	collectStore := func(l *runelexer.RuneLexer[WirToken]) {
		s := strings.TrimSpace(l.StoreFlush())
		if s != "" {
			l.TokenAppend(WirToken{
				t:    WirTokenTypeRawText,
				text: strings.TrimSpace(s),
			})
		}
	}
	ranFinal := false
	for {
		ch := l.Char()
		switch ch {
		default:
			{
				l.Store()
				break
			}
		case "@":
			{
				if l.InQuote() {
					l.Store()
					break
				}
				collectStore(l)
				if l.Pull(4) == "@for(" {
					l.Mark()
					l.NextUntil(")")
					l.TokenAppend(WirToken{
						t:    WirTokenTypeAtDirective,
						text: l.PullFromMark(),
					})
				} else {
					l.Store()
					break
				}
			}
		case "'":
			{
				collectStore(l)
				l.Mark()
				l.Next()
				l.Iter(func(ch2 string, pos int) bool {
					if ch2 == "'" && l.Peek(-1) != "\\" {
						l.TokenAppend(WirToken{
							t:    WirTokenTypeString,
							text: l.PullFromMark(),
						})
						return false
					}
					return true
				})
			}
		case "\"":
			{
				collectStore(l)
				l.Mark()
				l.Next()
				l.Iter(func(ch2 string, pos int) bool {
					if ch2 == "\"" && l.Peek(-1) != "\\" {
						l.TokenAppend(WirToken{
							t:    WirTokenTypeString,
							text: l.PullFromMark(),
						})
						return false
					}
					return true
				})
			}
		case "{":
			{
				if l.InQuote() {
					l.Store()
					break
				}
				collectStore(l)
				l.TokenAppend(WirToken{
					t:    WirTokenTypeCurlyBraceOpen,
					text: "{",
				})
			}
		case "}":
			{
				if l.InQuote() {
					l.Store()
					break
				}
				collectStore(l)
				l.TokenAppend(WirToken{
					t:    WirTokenTypeCurlyBraceClose,
					text: "}",
				})
			}
		case "<":
			{
				if l.InQuote() {
					l.Store()
					break
				}
				collectStore(l)
				if l.StoreLen() > 0 {
					l.TokenAppend(WirToken{
						t:    WirTokenTypeRawText,
						text: l.StoreFlush(),
					})
				}
				l.Mark()
				l.Iter(func(ch2 string, pos int) bool {
					if l.InQuote() {
						return true
					}
					if ch2 != ">" && ch2 != "}" {
						return true
					}
					if ch2 == "}" {
						l.TokenAppend(WirToken{
							t:    WirTokenTypeRawText,
							text: l.PullFromMark(),
						})
					}
					if ch2 == ">" {
						l.TokenAppend(WirToken{
							t:    WirTokenTypeTagInfo,
							text: l.PullFromMark(),
						})
					}
					return false
				})
			}
		}
		l.Next()
		if l.AtEnd() {
			if ranFinal {
				break
			} else {
				ranFinal = true
			}
		}
	}
	l.TokenAppend(WirToken{
		t:    WirTokenTypeEndOfFile,
		text: "EOF",
	})
	return nil
}
