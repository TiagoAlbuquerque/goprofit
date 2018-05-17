
package utils

import (
    "os"
    "os/exec"
    "strconv"
    "strings"
    "fmt"
    )

func terminalSize() (string, error) {
    cmd := exec.Command("stty", "size")
    cmd.Stdin = os.Stdin
    out, err := cmd.Output()
    return string(out), err
}

func parse(input string) (int, int, error) {
    parts := strings.Split(input, " ")
    x, err := strconv.Atoi(parts[0])
    if err != nil {
        return 0, 0, err
    }
    y, err := strconv.Atoi(strings.Replace(parts[1], "\n", "", 1))
    if err != nil {
        return 0, 0, err
    }
    return int(x), int(y), nil
}

    // Width return the width of the terminal.
func terminalWidth() (int, error) {
    output, err := terminalSize()
    if err != nil {
       return 0, err
    }
    _, width, err := parse(output)
    return width, err
}

    // Height returns the height of the terminal.
func terminalHeight() (int, error) {
    output, err := terminalSize()
    if err != nil {
       return 0, err
    }
    height, _, err := parse(output)
    return height, err
}

func ProgressBar(c int, t int){
    w, _ := terminalWidth()
    out := fmt.Sprintf("\r %d ", int((c*100)/t))
    out += "%% ["
    w -= len(out)+1
    out += strings.Repeat("#", int(c*w/t))
    w -= len(out)
    out += strings.Repeat(" ", w)
    out += "]"

    fmt.Printf(out)
}
