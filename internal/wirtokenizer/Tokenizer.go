package wirtokenizer

import (
	"os"
	"strings"

	"github.com/phillip-england/wir/internal/runelexer"
	"github.com/phillip-england/wir/internal/wherr"
)

type Tokenizer struct {
	Lexer *runelexer.RuneLexer[Token]
}

func TokenizerNewFromFile(path string) (*Tokenizer, error) {
	fBytes, err := os.ReadFile(path)
	if err != nil {
		return &Tokenizer{}, wherr.Consume(wherr.Here(), err, "")
	}
	fStr := string(fBytes)
	tk, err := TokenizerNewFromString(fStr)
	if err != nil {
		return &Tokenizer{}, wherr.Consume(wherr.Here(), err, "")
	}
	return tk, nil
}

func TokenizerNewFromString(s string) (*Tokenizer, error) {
	s = strings.TrimSpace(s)
	l := runelexer.NewRuneLexer[Token](s)
	err := tokenizeWir(l)
	if err != nil {
		return &Tokenizer{}, err
	}
	return &Tokenizer{
		Lexer: l,
	}, nil
}

func (t *Tokenizer) Str() string {
	s := ""
	if t.Lexer.TokenLen() == 0 {
		return ""
	}
	t.Lexer.TokenIter(func(token Token, index int) bool {
		s += string(token.t) + ":" + token.text + "\n"
		return true
	})
	return s
}

type TokenState = int

const (
	TokenStateInit = iota
)

func tokenizeWir(l *runelexer.RuneLexer[Token]) error {
	err := phase1(l)
	if err != nil {
		return err
	}
	err = phase2(l)
	if err != nil {
		return err
	}
	err = phase3(l)
	if err != nil {
		return err
	}
	return nil
}

func phase3(l *runelexer.RuneLexer[Token]) error {
	var toks []Token
	l.TokenIter(func(tk Token, index int) bool {
		switch tk.t {
		default:
			{
				toks = append(toks, tk)
			}
		case TokenTypeCurlyBraceOpen:
			{
				toks = append(toks, Token{
					t:    TokenTypeHTMLCurlyBraceOpen,
					text: tk.text,
				})
			}
		case TokenTypeCurlyBraceClose:
			{
				toks = append(toks, Token{
					t:    TokenTypeHTMLCurlyBraceClose,
					text: tk.text,
				})
			}
		case TokenTypeRawText:
			{
				toks = append(toks, Token{
					t:    TokenTypeHTMLTagName,
					text: tk.text,
				})
			}
		case TokenTypeDollarSignInterpolation:
			{
				toks = append(toks, Token{
					t:    TokenTypeDollarSignInterpolationOpen,
					text: "${",
				})
				s := tk.text
				s = strings.Replace(s, "${", "", 1)
				s = s[0 : len(s)-1]
				innerParts := strings.Split(s, ":")
				for i, v := range innerParts {
					if i%2 == 0 {
						toks = append(toks, Token{
							t:    TokenTypeDollarSignInterpolationValue,
							text: strings.TrimSpace(v),
						})
						toks = append(toks, Token{
							t:    TokenTypeDollarSignInterpolationSemiColon,
							text: ":",
						})
					} else {
						toks = append(toks, Token{
							t:    TokenTypeDollarSignInterpolationType,
							text: strings.TrimSpace(v),
						})
					}
				}
			}
			toks = append(toks, Token{
				t:    TokenTypeDollarSignInterpolationClose,
				text: "}",
			})
		}
		return true
	})
	l.TokenOverwrite(toks)
	return nil
}

func phase2(l *runelexer.RuneLexer[Token]) error {
	var toks []Token
	l.TokenIter(func(tk Token, index int) bool {
		switch tk.t {
		default:
			{
				toks = append(toks, tk)
			}
		case TokenTypeTagInfo:
			{
				l2 := runelexer.NewRuneLexer[Token](tk.text)
				l2.Iter(func(ch string, pos int) bool {
					switch ch {
					default:
						{
							l2.Store()
						}
					case "=":
						{
							attrKey := strings.TrimSpace(l2.StoreFlush())
							toks = append(toks, Token{
								t:    TokenTypeHTMLAttrKey,
								text: attrKey,
							})
							toks = append(toks, Token{
								t:    TokenTypeHTMLAttrEqualSign,
								text: "=",
							})
						}
					case "<":
						{
							toks = append(toks, Token{
								t:    TokenTypeHTMLTagInfoStart,
								text: "<",
							})
						}
					case ">":
						{
							attrKey := strings.TrimSpace(l2.StoreFlush())
							if attrKey != "" {
								toks = append(toks, Token{
									t:    TokenTypeHTMLAttrKey,
									text: attrKey,
								})
							}
							toks = append(toks, Token{
								t:    TokenTypeHTMLTagInfoEnd,
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
									l3 := runelexer.NewRuneLexer[Token](htmlAttr)
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
											toks = append(toks, Token{
												t:    TokenTypeHTMLAttrValuePartial,
												text: attrBit,
											})
											l3.Mark()
											l3.NextUntil("}")
											dollarSignInterpolation := l3.PullFromMark()
											toks = append(toks, Token{
												t:    TokenTypeDollarSignInterpolation,
												text: dollarSignInterpolation,
											})
										}
										return true
									})
									if brokeAttr {
										toks = append(toks, Token{
											t:    TokenTypeHTMLAttrValuePartial,
											text: l3.StoreFlush(),
										})
									} else {
										toks = append(toks, Token{
											t:    TokenTypeHTMLAttrValue,
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
									toks = append(toks, Token{
										t:    TokenTypeHTMLAttrValue,
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
		case TokenTypeString:
			{
				l2 := runelexer.NewRuneLexer[Token](tk.text)
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
							toks = append(toks, Token{
								t:    TokenTypeStringContent,
								text: l2.StoreFlush(),
							})
							l2.Mark()
							l2.NextUntil("}")
							toks = append(toks, Token{
								t:    TokenTypeDollarSignInterpolation,
								text: l2.PullFromMark(),
							})
						}
					case "'":
						{
							if pos == 0 {
								toks = append(toks, Token{
									t:    TokenTypeStringStart,
									text: "'",
								})
								return true
							}
							if l2.AtEnd() && ch == "'" {
								flush := l2.StoreFlush()
								if flush != "" {
									toks = append(toks, Token{
										t:    TokenTypeStringContent,
										text: flush,
									})
								}
								toks = append(toks, Token{
									t:    TokenTypeStringEnd,
									text: "'",
								})
								return true
							}
						}
					case "\"":
						{
							if pos == 0 {
								toks = append(toks, Token{
									t:    TokenTypeStringStart,
									text: "\"",
								})
								return true
							}
							if l2.AtEnd() && ch == "\"" {
								toks = append(toks, Token{
									t:    TokenTypeStringContent,
									text: l2.StoreFlush(),
								})
								toks = append(toks, Token{
									t:    TokenTypeStringEnd,
									text: "\"",
								})
								return true
							}
						}
					}
					return true
				})
			}
		case TokenTypeAtDirective:
			{
				l2 := runelexer.NewRuneLexer[Token](tk.text)
				l2.Iter(func(ch string, pos int) bool {
					switch ch {
					default:
						{
							l2.Store()
						}
					case ")":
						{
							toks = append(toks, Token{
								t:    TokenTypeAtDirectiveParenthesisClose,
								text: ")",
							})
						}
					case "(":
						{
							directiveName := l2.StoreFlush()
							toks = append(toks, Token{
								t:    TokenTypeAtDirectiveName,
								text: directiveName,
							})
							toks = append(toks, Token{
								t:    TokenTypeAtDirectiveParenthesisOpen,
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
									toks = append(toks, Token{
										t:    TokenTypeAtDirectiveParamValue,
										text: v,
									})
									toks = append(toks, Token{
										t:    TokenTypeAtDirectiveSemiColon,
										text: ":",
									})
								} else {
									toks = append(toks, Token{
										t:    TokenTypeAtDirectiveParamType,
										text: v,
									})
								}
							}
						}
					case "@":
						{
							toks = append(toks, Token{
								t:    TokenTypeAtDirectiveStart,
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

func phase1(l *runelexer.RuneLexer[Token]) error {
	collectStore := func(l *runelexer.RuneLexer[Token]) {
		flush := l.StoreFlush()
		s := strings.TrimSpace(flush)
		if s != "" {
			l.TokenAppend(Token{
				t:    TokenTypeRawText,
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
					l.TokenAppend(Token{
						t:    TokenTypeAtDirective,
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
						l.TokenAppend(Token{
							t:    TokenTypeString,
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
						l.TokenAppend(Token{
							t:    TokenTypeString,
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
				l.TokenAppend(Token{
					t:    TokenTypeCurlyBraceOpen,
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
				l.TokenAppend(Token{
					t:    TokenTypeCurlyBraceClose,
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
					l.TokenAppend(Token{
						t:    TokenTypeRawText,
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
						l.TokenAppend(Token{
							t:    TokenTypeRawText,
							text: l.PullFromMark(),
						})
					}
					if ch2 == ">" {
						l.TokenAppend(Token{
							t:    TokenTypeTagInfo,
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
	l.TokenAppend(Token{
		t:    TokenTypeEndOfFile,
		text: "EOF",
	})
	return nil
}