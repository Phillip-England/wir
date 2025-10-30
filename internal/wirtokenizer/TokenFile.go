package wirtokenizer

import (
	"os"

	"github.com/phillip-england/wir/internal/wherr"
)

type TokenFile struct {
	Path string
	Text string
	Tokens []Token
	TokenStr string
}

func TokenFileNewFromTokenizer(path string, tk *Tokenizer) (TokenFile, error) {
	s := tk.Str()
	toks, err := TokenManyFromStr(s)
	if err != nil {
		return TokenFile{}, wherr.Consume(wherr.Here(), err, "")
	}
	tokenStr := ""
	for _, tok := range toks {
		tokenStr += tok.text
	}
	return TokenFile{
		Path: path,
		Text: s,
		Tokens: toks,
		TokenStr: tokenStr,
	}, nil
}

func TokenFileLoad(path string) (TokenFile, error) {
	fBytes, err := os.ReadFile(path)
	if err != nil {
		return TokenFile{}, wherr.Consume(wherr.Here(), err, "")
	}
	fStr := string(fBytes)
	toks, err := TokenManyFromStr(fStr)
	if err != nil {
		return TokenFile{}, wherr.Consume(wherr.Here(), err, "")
	}
	tokenStr := ""
	for _, tok := range toks {
		tokenStr += tok.text
	}
	return TokenFile{
		Path: path,
		Text: fStr,
		Tokens: toks,
		TokenStr: tokenStr,
	}, nil
}


