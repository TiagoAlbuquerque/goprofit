// Package color provides functions for formatting text in different colors.
package color

import "fmt"

// ANSI escape code for resetting the text color to the default.
const reset = "\033[0m"

// ANSI escape code for setting the text color using an 8-bit color value.
const c256 = "\033[38;5;%dm"

// ANSI escape code for setting the text color using an RGB color value.
const cRGB = "\033[38;2;%d;%d;%dm"

// Fg8b formats the string foreground (text) to output in the specified 8-bit color.
// It takes an integer c representing the 8-bit color value (0-255) and a string txt as input.
// It returns the formatted string with the specified color.
func Fg8b(c int, txt string) string {
	return fmt.Sprintf(c256+txt+reset, c)
}
