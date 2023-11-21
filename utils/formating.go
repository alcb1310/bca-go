package utils

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func FormatNumber(f float64) string {
	p := message.NewPrinter(language.Spanish)
	return p.Sprintf("%.2f", f)
}
