package glam

import "fmt"

type ColorCode = string

const (
	ColorCodeRed     = "31"
	ColorCodeGreen   = "32"
	ColorCodeYellow  = "33"
	ColorCodeBlue    = "34"
	ColorCodeMagenta = "35"
	ColorCodeCyan    = "36"
	ColorCodeWhite   = "37"
)

// Wrap wraps text with ANSI escape codes for the given color.
func Wrap(code ColorCode, text string) string {
	return fmt.Sprintf("\033[%sm%s\033[0m", code, text)
}

// Color functions — behave like fmt.Sprintf but return colored text
func Red(format string, args ...any) string {
	return Wrap(ColorCodeRed, fmt.Sprintf(format, args...))
}

func Green(format string, args ...any) string {
	return Wrap(ColorCodeGreen, fmt.Sprintf(format, args...))
}

func Yellow(format string, args ...any) string {
	return Wrap(ColorCodeYellow, fmt.Sprintf(format, args...))
}

func Blue(format string, args ...any) string {
	return Wrap(ColorCodeBlue, fmt.Sprintf(format, args...))
}

func Magenta(format string, args ...any) string {
	return Wrap(ColorCodeMagenta, fmt.Sprintf(format, args...))
}

func Cyan(format string, args ...any) string {
	return Wrap(ColorCodeCyan, fmt.Sprintf(format, args...))
}

func White(format string, args ...any) string {
	return Wrap(ColorCodeWhite, fmt.Sprintf(format, args...))
}
