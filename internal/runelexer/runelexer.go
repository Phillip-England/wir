package runelexer

import (
	"fmt"
	"math"
	"runtime"
)

type RuneLexer[T any] struct {
	position    int
	tokPosition int
	endPosition int
	runes       []rune
	markedPos   int
	tokens      []T
	state       int
	store       []rune
}

func NewRuneLexer[T any](s string) *RuneLexer[T] {
	runes := []rune(s)
	return &RuneLexer[T]{
		position:    0,
		tokPosition: 0,
		endPosition: len(runes) - 1,
		runes:       runes,
		markedPos:   0,
		tokens:      []T{},
		state:       0,
	}
}

func (l *RuneLexer[T]) State() int {
	return l.state
}

func (l *RuneLexer[T]) SetState(n int) {
	l.state = n
}

func (l *RuneLexer[T]) Runes() []rune {
	return l.runes
}

func (l *RuneLexer[T]) Char() string {
	if len(l.runes) == 0 {
		return ""
	}
	return string(l.runes[l.position])
}

func (l *RuneLexer[T]) AtEnd() bool {
	return l.position >= l.endPosition
}

func (l *RuneLexer[T]) AtStart() bool {
	return l.position == 0
}

func (l *RuneLexer[T]) Len() int {
	return len(l.runes)
}

func (l *RuneLexer[T]) Pos() int {
	return l.position
}

func (l *RuneLexer[T]) Next() {
	l.position++
	if l.AtEnd() {
		l.position = l.endPosition
	}
}

func (l *RuneLexer[T]) Prev() {
	l.position--
	if l.position < 0 {
		l.position = 0
	}
}

func (l *RuneLexer[T]) Iter(fn func(ch string, pos int) bool) {
	for {
		shouldContinue := fn(l.Char(), l.Pos())
		if !shouldContinue {
			break
		}
		if l.AtEnd() {
			break
		}
		l.Next()
	}
}

func (l *RuneLexer[T]) Mark() {
	l.markedPos = l.Pos()
}

func (l *RuneLexer[T]) GoToEnd() {
	l.position = l.endPosition
}

func (l *RuneLexer[T]) GoToStart() {
	l.position = 0
}

func (l *RuneLexer[T]) GoToMark() {
	l.position = l.markedPos
}

func (l *RuneLexer[T]) PullFromStart() string {
	start := 0
	end := l.position + 1
	if end > len(l.runes) {
		end = len(l.runes)
	}
	return string(l.runes[start:end])
}

func (l *RuneLexer[T]) PullFromEnd() string {
	start := l.position
	if start < 0 {
		start = 0
	}
	return string(l.runes[start:])
}

func (l *RuneLexer[T]) PullFromMark() string {
	start := l.markedPos
	end := l.position
	if start > end {
		start, end = end, start
	}
	end++
	if end > len(l.runes) {
		end = len(l.runes)
	}
	return string(l.runes[start:end])
}

func (l *RuneLexer[T]) NextBy(n int) {
	for count := 0; count < n; count++ {
		l.Next()
	}
}

func (l *RuneLexer[T]) PrevBy(n int) {
	for count := 0; count < n; count++ {
		l.Prev()
	}
}

func (l *RuneLexer[T]) Peek(n int) string {
	start := l.position
	if n >= 0 {
		l.NextBy(n)
	} else {
		l.PrevBy(int(math.Abs(float64(n))))
	}
	c := l.Char()
	l.position = start
	return c
}

func (l *RuneLexer[T]) Str() string {
	return string(l.runes)
}

func (l *RuneLexer[T]) NextUntil(ch string) bool {
	l.Next()
	for !l.AtEnd() {
		if l.Char() == ch {
			return true
		}
		l.Next()
	}
	return l.Char() == ch
}

func (l *RuneLexer[T]) PrevUntil(ch string) bool {
	l.Prev()
	for l.position > 0 {
		if l.Char() == ch {
			return true
		}
		l.Prev()
	}
	return l.Char() == ch
}

func (l *RuneLexer[T]) NextUntilNot(ch string) bool {
	if l.Char() != ch {
		return true
	}
	for !l.AtEnd() {
		l.Next()
		if l.Char() != ch {
			return true
		}
	}
	return l.Char() != ch
}

func (l *RuneLexer[T]) PrevUntilNot(ch string) bool {
	if l.Char() != ch {
		return true
	}
	for l.position > 0 {
		l.Prev()
		if l.Char() != ch {
			return true
		}
	}
	return l.Char() != ch
}

func (l *RuneLexer[T]) PullFrom(start, end int) string {
	if start < 0 {
		start = 0
	}
	if end < 0 {
		end = 0
	}
	if start >= len(l.runes) {
		start = len(l.runes) - 1
	}
	if end >= len(l.runes) {
		end = len(l.runes) - 1
	}
	if start > end {
		start, end = end, start
	}
	end++
	if end > len(l.runes) {
		end = len(l.runes)
	}
	return string(l.runes[start:end])
}

func (l *RuneLexer[T]) MarkedPos() int {
	return l.markedPos
}

func (l *RuneLexer[T]) NextUntilAny(chars ...string) string {
	start := l.position
	for !l.AtEnd() {
		ch := string(l.runes[l.position])
		if contains(chars, ch) {
			break
		}
		l.position++
	}
	return string(l.runes[start:l.position])
}

func (l *RuneLexer[T]) PrevUntilAny(chars ...string) string {
	start := l.position
	for !l.AtStart() {
		ch := string(l.runes[l.position])
		if contains(chars, ch) {
			break
		}
		l.position--
	}
	return string(l.runes[l.position+1 : start+1])
}

func (l *RuneLexer[T]) NextUntilNotAny(chars ...string) string {
	start := l.position
	for !l.AtEnd() {
		ch := string(l.runes[l.position])
		if !contains(chars, ch) {
			break
		}
		l.position++
	}
	return string(l.runes[start:l.position])
}

func (l *RuneLexer[T]) PrevUntilNotAny(chars ...string) string {
	start := l.position
	for !l.AtStart() {
		ch := string(l.runes[l.position])
		if !contains(chars, ch) {
			break
		}
		l.position--
	}
	return string(l.runes[l.position+1 : start+1])
}

func (l *RuneLexer[T]) Here() location {
	return here()
}

func (l *RuneLexer[T]) InQuote() bool {
	return isInQuote(l.runes, l.position)
}

func (l *RuneLexer[T]) TokenAppend(token T) {
	l.tokens = append(l.tokens, token)
}

func (l *RuneLexer[T]) TokenIter(fn func(tk T, index int) bool) {
	for i, token := range l.tokens {
		shouldContinue := fn(token, i)
		if !shouldContinue {
			break
		}
	}
}

func (l *RuneLexer[T]) TokenLast() T {
	var zero T
	if len(l.tokens) == 0 {
		return zero
	}
	return l.tokens[len(l.tokens)-1]
}

func (l *RuneLexer[T]) TokenLen() int {
	return len(l.tokens)
}

func (l *RuneLexer[T]) Tokens() []T {
	return l.tokens
}

func (l *RuneLexer[T]) TokenClean() {
	l.tokens = []T{}
}

func (l *RuneLexer[T]) TokenOverwrite(toks []T) {
	l.tokens = toks
}

func contains(arr []string, target string) bool {
	for _, a := range arr {
		if a == target {
			return true
		}
	}
	return false
}

func (l *RuneLexer[T]) Store() {
	if l.position >= 0 && l.position < len(l.runes) {
		l.store = append(l.store, l.runes[l.position])
	}
}

func (l *RuneLexer[T]) StoreStr() string {
	return string(l.store)
}

func (l *RuneLexer[T]) StoreClear() {
	l.store = []rune{}
}

func (l *RuneLexer[T]) StoreLen() int {
	return len(l.store)
}

func (l *RuneLexer[T]) StoreFlush() string {
	s := l.StoreStr()
	l.StoreClear()
	return s
}

func (l *RuneLexer[T]) Pull(n int) string {
	if n == 0 {
		return l.Char()
	}

	var start, end int

	if n > 0 {
		start = l.position
		end = l.position + n + 1
		if end > len(l.runes) {
			end = len(l.runes)
		}
	} else {
		start = l.position + n
		end = l.position + 1
		if start < 0 {
			start = 0
		}
		if end > len(l.runes) {
			end = len(l.runes)
		}
	}

	return string(l.runes[start:end])
}

type location struct {
	File string
	Line int
}

func (l location) str() string {
	return fmt.Sprintf("Line: %d: File: %s", l.Line, l.File)
}

func here() location {
	_, file, line, _ := runtime.Caller(1)
	return location{File: file, Line: line}
}

func isInQuote(r []rune, pos int) bool {
	inDouble := false
	inSingle := false
	inQuote := false
	for i, rn := range r {

		prevChar := ""
		if i > 0 {
			prevChar = string(r[i-1])
		}
		if i > pos {
			break
		}
		ch := string(rn)
		switch ch {
		default:
			{
				continue
			}
		case "\"":
			{
				if prevChar == "\\" {
					continue
				}
				if !inDouble && !inSingle {
					inDouble = true
					inQuote = true
					continue
				}
				if inDouble && !inSingle {
					inDouble = false
					inQuote = false
					continue
				}
				if inDouble && inSingle {
					inDouble = false
					continue
				}
				if !inDouble && inSingle {
					inQuote = true
					continue
				}
			}
		case "'":
			{
				if prevChar == "\\" {
					continue
				}
				if !inSingle && !inDouble {
					inSingle = true
					inQuote = true
					continue
				}
				if inSingle && !inDouble {
					inSingle = false
					inQuote = false
					continue
				}
				if inSingle && inDouble {
					inSingle = false
					continue
				}
				if !inSingle && inDouble {
					inQuote = true
					continue
				}
			}
		}
	}
	return inQuote
}
