package utils

import (
	"fmt"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

func colorize(color, s string) string {
	return color + s + colorReset
}

func Info(s string, args ...interface{}) {
	fmt.Println(colorize(colorCyan, fmt.Sprintf(s, args...)))
}

func Success(s string, args ...interface{}) {
	fmt.Println(colorize(colorGreen, fmt.Sprintf(s, args...)))
}

func Error(s string, args ...interface{}) {
	fmt.Println(colorize(colorRed, fmt.Sprintf(s, args...)))
}

func Prompt(s string, args ...interface{}) {
	fmt.Print(colorize(colorYellow, fmt.Sprintf(s, args...)))
}

func Note(s string, args ...interface{}) {
	fmt.Println(colorize(colorBlue, fmt.Sprintf(s, args...)))
}

func Warn(s string, args ...interface{}) {
	fmt.Println(colorize(colorYellow, fmt.Sprintf("[WARN] "+s, args...)))
}
