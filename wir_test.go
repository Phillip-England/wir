package webir_test

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/phillip-england/wir/internal/soak"
	"github.com/phillip-england/wir/internal/wherr"
	"github.com/phillip-england/wir/internal/wirtokenizer"
)

func fail(t *testing.T, err error) {
	fmt.Println(err.Error())
	t.Fail()
}

func TestWirParser(t *testing.T) {
	fmt.Println("testing")
}

func TestLoadTokenFile(t *testing.T) {
	cwd, _ := os.Getwd()
	_, err := wirtokenizer.TokenFileLoad(path.Join(cwd, "examples", "toks", "user_list.tok"))
	if err != nil {
		fail(t, wherr.Consume(wherr.Here(), err, ""))
	}
}

func TestExamplesTokenized(t *testing.T) {
	d, err := soak.NewMirror([]string{"examples", "raw"}, []string{"examples", "toks"})
	if err != nil {
		fail(t, wherr.Consume(wherr.Here(), err, ""))
	}
	d.Iter(func(target *soak.VirtualAsset, compare *soak.VirtualAsset) bool {
		compareTokFile, err := wirtokenizer.TokenFileLoad(compare.Path)
		if err != nil {
			fail(t, wherr.Consume(wherr.Here(), err, ""))
		}
		tk, err := wirtokenizer.TokenizerNewFromString(target.Text)
		if err != nil {
			fail(t, wherr.Consume(wherr.Here(), err, ""))
		}
		targetTokFile, err := wirtokenizer.TokenFileNewFromTokenizer("", tk)
		if err != nil {
			fail(t, wherr.Consume(wherr.Here(), err, ""))
		}
		targetStr := targetTokFile.TokenStr
		compareStr := compareTokFile.TokenStr
		if len(targetStr) != len(compareStr) {
			fail(t, wherr.Err(wherr.Here(), "target file: [%s] does not mirror comparison file: [%s]", target.Path, compare.Path))
		}
		return true
	})
	if err != nil {
		fail(t, wherr.Consume(wherr.Here(), err, ""))
	}
}
