package glam

import "fmt"

type ColorCode = string

const (
	ColorCodeRed     = "RED"
	ColorCodeGreen   = "GREEN"
	ColorCodeYellow  = "YELLOW"
	ColorCodeBlue    = "BLUE"
	ColorCodeMagenta = "MAGENTA"
	ColorCodeCyan    = "CYAN"
	ColorCodeWhite   = "WHITE"
)

type Glam struct {
	ColorCode ColorCode
	Text      string
}

func New(code ColorCode, text string) Glam {
	return Glam{
		ColorCode: code,
		Text:      text,
	}
}

// Color functions
func Red(format string, args ...any) Glam {
	return New(ColorCodeRed, fmt.Sprintf(format, args...))
}

func Green(format string, args ...any) Glam {
	return New(ColorCodeGreen, fmt.Sprintf(format, args...))
}

func Yellow(format string, args ...any) Glam {
	return New(ColorCodeYellow, fmt.Sprintf(format, args...))
}

func Blue(format string, args ...any) Glam {
	return New(ColorCodeBlue, fmt.Sprintf(format, args...))
}

func Magenta(format string, args ...any) Glam {
	return New(ColorCodeMagenta, fmt.Sprintf(format, args...))
}

func Cyan(format string, args ...any) Glam {
	return New(ColorCodeCyan, fmt.Sprintf(format, args...))
}

func White(format string, args ...any) Glam {
	return New(ColorCodeWhite, fmt.Sprintf(format, args...))
}

// Wrap wraps the input text with the proper ANSI color code
func Wrap(code ColorCode, text string) string {
	ansiCode := getANSICode(code)
	return fmt.Sprintf("\033[%sm%s\033[0m", ansiCode, text)
}

// Print takes Glam instances and prints them with colors
func Print(glams ...Glam) {
	for _, g := range glams {
		fmt.Print(Wrap(g.ColorCode, g.Text))
	}
}

// getANSICode returns the ANSI color code for a given ColorCode
func getANSICode(code ColorCode) string {
	switch code {
	case ColorCodeRed:
		return "31"
	case ColorCodeGreen:
		return "32"
	case ColorCodeYellow:
		return "33"
	case ColorCodeBlue:
		return "34"
	case ColorCodeMagenta:
		return "35"
	case ColorCodeCyan:
		return "36"
	case ColorCodeWhite:
		return "37"
	default:
		return "0" // Reset/default
	}
}
