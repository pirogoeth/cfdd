package util

import (
	"fmt"

	"github.com/mgutz/ansi"
)

var (
	Okay  func(s string) string = ansi.ColorFunc("green+h")
	Warn  func(s string) string = ansi.ColorFunc("yellow")
	Error func(s string) string = ansi.ColorFunc("red")
	Fatal func(s string) string = ansi.ColorFunc("red+uh")
)

func Okayf(msg string, args ...interface{}) string {
	return Okay(fmt.Sprintf(msg, args...))
}

func Warnf(msg string, args ...interface{}) string {
	return Warn(fmt.Sprintf(msg, args...))
}

func Errorf(msg string, args ...interface{}) string {
	return Error(fmt.Sprintf(msg, args...))
}

func Fatalf(msg string, args ...interface{}) string {
	return Fatal(fmt.Sprintf(msg, args...))
}
