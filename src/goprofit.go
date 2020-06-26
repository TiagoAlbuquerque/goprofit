package main

import (
	"./controller"
	"fmt"
	"os"
	"os/signal"
)

func main() {
	for {
		cicle()
	}
}

func cicle() {
	controller.FetchMarket()
	controller.PrintShoppingLists(3)

	controller.Terminate()
}

func interruptionHandler(c chan os.Signal) {
	for sig := range c {
		fmt.Println(sig)
		controller.Terminate()
		os.Exit(0)
	}
}
func init() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go interruptionHandler(c)
}
