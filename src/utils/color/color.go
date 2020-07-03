package color

import "fmt"

const reset = "\033[0m"
const c256 = "\033[38;5;%dm"
const cRGB = "\033[38;2;%d;%d;%dm"

//Fg8b will format the string foreground (text) to output in the specified 8-bit color
func Fg8b(c int, txt string) string {
	return fmt.Sprintf(c256+txt+reset, c)
}

//FgRGB will format the string foreground (text) to output in the specified rgb color
func FgRGB(r, g, b byte, txt string) string {
	return fmt.Sprintf(cRGB+txt+reset, r, g, b)
}
