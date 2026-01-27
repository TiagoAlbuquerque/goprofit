package main

import (
	"fmt"
	"os"
	"os/signal"

	"goprofit/controller"
	"goprofit/server"
	shoppinglists "goprofit/shoppingLists"
)

func main() {
	go server.Start()

	for {
		cicle()
	}
}

func cicle() {
	shoppinglists.NextRound()
	controller.FetchMarket()
	shoppinglists.Prune()
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
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go interruptionHandler(c)
}
