package webir_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/phillip-england/wir/internal/glam"
	"github.com/phillip-england/wir/internal/runelexer"
	"github.com/phillip-england/wir/internal/soak"
	"github.com/phillip-england/wir/internal/wherr"
	"github.com/phillip-england/wir/internal/wirtokenizer"
)

func fail(t *testing.T, err error) {
	fmt.Println(err.Error())
	t.Fail()
}

func TestExamplesTokenized(t *testing.T) {
	d, err := soak.NewMirror([]string{"examples", "raw"}, []string{"examples", "tokenized"})
	if err != nil {
		fail(t, wherr.Consume(wherr.Here(), err, ""))
	}
	d.Iter(func(target *soak.VirtualAsset, compare *soak.VirtualAsset) bool {
		tk, err := wirtokenizer.NewFromString(target.Text)
		if err != nil {
			fail(t, wherr.Consume(wherr.Here(), err, ""))
		}
		tokStr := tk.TokStr()
		if tokStr != compare.Text {
			l1 := runelexer.NewRuneLexer[any](tokStr)
			l2 := runelexer.NewRuneLexer[any](compare.Text)
			badPosition := 0
			l1.Iter(func(ch1 string, pos1 int) bool {
				if pos1 > len(l2.Runes())-1 || pos1 < 0 {
					return true
				}
				ch2 := string(l2.Runes()[pos1])
				if ch1 != ch2 {
					if badPosition == len(l2.Runes())-1 {
						badPosition = pos1
					} else {
						badPosition = pos1 + 1
					}
					return false
				}
				return true
			})
			lineNumber := strings.Count(tokStr[0:badPosition], "\n")
			targetLines := strings.Split(tokStr, "\n")
			compareLines := strings.Split(compare.Text, "\n")
			targetLine := targetLines[lineNumber]
			compareLine := compareLines[lineNumber]
			glam.Print(
				glam.Red("[DIR COMPARE FAILURE]\n"),
				glam.Yellow("[TARGET] => %s\n", target.RelPath),
				glam.Yellow("[COMPARED TO] => %s\n", compare.RelPath),
				glam.Yellow("[LINE] => %d\n", lineNumber+1),
				glam.Green("[TARGET] => %s\n", targetLine),
				glam.Red("[ACTUAL] => %s\n", compareLine),
			)
			t.Fail()
		}
		return true
	})
	if err != nil {
		fail(t, wherr.Consume(wherr.Here(), err, ""))
	}
}
