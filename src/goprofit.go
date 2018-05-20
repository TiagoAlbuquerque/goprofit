package main

import(
    "./regions"
    "./items"
    "fmt"
    "os"
    "os/signal"
)

func main(){
    defer func() {
        recover()
        terminate()
    }()

    cicle()
}
func cicle() {
    for {
        regions.GetMarketsPages()
    }
}
func terminate(){
    fmt.Println("terminating")
    regions.Terminate()
    items.Terminate()
}

func interruptionHandler(c chan os.Signal){
    for sig := range c {
        fmt.Println(sig)
        terminate()
        os.Exit(0)
    }
}
func init(){
    c := make(chan os.Signal)
    signal.Notify(c, os.Interrupt)
    go interruptionHandler(c)
}
