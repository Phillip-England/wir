package runelexer

import (
	"math"
)

type AbstractLexer[T any] struct {
	position       int
	endPosition    int
	markedPosition int
	items          []T 
	store          []T
}

func AbstractLexerNew[T any](items []T) *AbstractLexer[T] {
	return &AbstractLexer[T]{
		position:       0,
		endPosition:    len(items) - 1, 
		markedPosition: 0,
		items:          items, 
		store:          []T{},
	}
}

func (l *AbstractLexer[T]) Items() []T {
	return l.items
}

func (l *AbstractLexer[T]) Item() T {
	var zero T
	if len(l.items) == 0 {
		return zero
	}
	if l.position < 0 || l.position >= len(l.items) {
		return zero
	}
	return l.items[l.position]
}

func (l *AbstractLexer[T]) AtEnd() bool {
	return l.position >= l.endPosition
}

func (l *AbstractLexer[T]) AtStart() bool {
	return l.position == 0
}

func (l *AbstractLexer[T]) Len() int {
	return len(l.items)
}

func (l *AbstractLexer[T]) Pos() int {
	return l.position
}

func (l *AbstractLexer[T]) Next() {
	l.position++
	if l.AtEnd() {
		l.position = l.endPosition
	}
}

func (l *AbstractLexer[T]) Prev() {
	l.position--
	if l.position < 0 {
		l.position = 0
	}
}

func (l *AbstractLexer[T]) Iter(fn func(item T, pos int) bool) {
	if len(l.items) == 0 { 
		return
	}
	for {
		shouldContinue := fn(l.Item(), l.Pos())
		if !shouldContinue {
			break
		}
		if l.AtEnd() {
			break
		}
		l.Next()
	}
}

func (l *AbstractLexer[T]) Mark() {
	l.markedPosition = l.Pos()
}

func (l *AbstractLexer[T]) GoToEnd() {
	l.position = l.endPosition
}

func (l *AbstractLexer[T]) GoToStart() {
	l.position = 0
}

func (l *AbstractLexer[T]) GoToMark() {
	l.position = l.markedPosition
}

func (l *AbstractLexer[T]) PullFromStart() []T {
	start := 0
	end := l.position + 1
	if end > len(l.items) {
		end = len(l.items)
	}
	if start >= end {
		return []T{}
	}
	return l.items[start:end]
}

func (l *AbstractLexer[T]) PullFromEnd() []T {
	start := l.position
	if start < 0 {
		start = 0
	}
	if start >= len(l.items) {
		return []T{}
	}
	return l.items[start:]
}

func (l *AbstractLexer[T]) PullFromMark() []T {
	start := l.markedPosition
	end := l.position
	if start > end {
		start, end = end, start
	}
	end++
	if end > len(l.items) {
		end = len(l.items)
	}
	if start < 0 {
		start = 0
	}
	if start >= end {
		return []T{}
	}
	return l.items[start:end]
}

func (l *AbstractLexer[T]) NextBy(n int) {
	for count := 0; count < n; count++ {
		l.Next()
	}
}

func (l *AbstractLexer[T]) PrevBy(n int) {
	for count := 0; count < n; count++ {
		l.Prev()
	}
}

func (l *AbstractLexer[T]) Peek(n int) T {
	start := l.position
	if n >= 0 {
		l.NextBy(n)
	} else {
		l.PrevBy(int(math.Abs(float64(n))))
	}
	item := l.Item()
	l.position = start
	return item
}

func (l *AbstractLexer[T]) PullFrom(start, end int) []T {
	if start < 0 {
		start = 0
	}
	if end < 0 {
		end = 0
	}
	if start >= len(l.items) {
		start = len(l.items) - 1
	}
	if end >= len(l.items) {
		end = len(l.items) - 1
	}
	if start < 0 { 
		start = 0
	}
	if end < 0 {
		end = 0
	}
	if start > end {
		start, end = end, start
	}
	end++
	if end > len(l.items) {
		end = len(l.items)
	}
	if start >= end {
		return []T{}
	}
	return l.items[start:end]
}

func (l *AbstractLexer[T]) MarkedPos() int {
	return l.markedPosition
}

func (l *AbstractLexer[T]) Store() {
	if l.position >= 0 && l.position < len(l.items) {
		l.store = append(l.store, l.items[l.position])
	}
}

func (l *AbstractLexer[T]) StoreItems() []T {
	return l.store
}

func (l *AbstractLexer[T]) StoreClear() {
	l.store = []T{}
}

func (l *AbstractLexer[T]) StoreLen() int {
	return len(l.store)
}

func (l *AbstractLexer[T]) StoreFlush() []T {
	s := l.StoreItems()
	l.StoreClear()
	return s
}

func (l *AbstractLexer[T]) Pull(n int) []T {
	if len(l.items) == 0 {
		return []T{}
	}

	if n == 0 {
		return []T{l.Item()}
	}

	var start, end int

	if n > 0 {
		start = l.position
		end = l.position + n + 1 
		if end > len(l.items) {
			end = len(l.items)
		}
	} else {
		start = l.position + n   
		end = l.position + 1 
		if start < 0 {
			start = 0
		}
		if end > len(l.items) {
			end = len(l.items)
		}
	}
	if start < 0 {
		start = 0
	}
	if end > len(l.items) {
		end = len(l.items)
	}
	if start >= len(l.items) {
		return []T{}
	}
	if start >= end {
		return []T{}
	}

	return l.items[start:end]
}