package color

var colors = [...]string{
	"\033[0m",  //reset 0
	"\033[31m", //red 1
	"\033[32m", //green 2
	"\033[33m", //yellow 3
	"\033[34m", //blue 4
	"\033[35m", //purple 5
	"\033[36m", //cyan 6
	"\033[37m", //gray 7
}

//Fg will format the string foreground (text) to output in the specified color
func Fg(color int, txt string) string {
	if color >= len(colors) {
		return txt
	}
	return colors[color] + txt + colors[0]
}
