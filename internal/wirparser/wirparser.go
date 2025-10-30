package wirparser

import (
	"fmt"

	"github.com/phillip-england/wir/internal/runelexer"
	"github.com/phillip-england/wir/internal/wherr"
	"github.com/phillip-england/wir/internal/wirtokenizer"
)

type Parser struct {
	lexer *runelexer.AbstractLexer[wirtokenizer.Token]
}

func ParserNew(toks []wirtokenizer.Token) (*Parser, error) {
	l := runelexer.AbstractLexerNew(toks)
	l.Iter(func(item wirtokenizer.Token, pos int) bool {
		fmt.Println(item.Str())
		return true
	})
	err := recursiveParse(l)
	if err != nil {
		return &Parser{}, wherr.Consume(wherr.Here(), err, "")
	}
	return &Parser{
		lexer: l,
	}, nil
}

func recursiveParse(l *runelexer.AbstractLexer[wirtokenizer.Token]) (error) {
	l.Iter(func(item wirtokenizer.Token, pos int) bool {
		return true
	})	
	return nil
}

