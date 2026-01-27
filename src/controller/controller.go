package controller

import (
	"goprofit/conf"
	"goprofit/deals"
	"goprofit/items"
	"goprofit/locations"
	"goprofit/orders"
	"goprofit/regions"
	shoppinglists "goprofit/shoppingLists"
	"goprofit/utils"

	"encoding/json"
	"fmt"
	"io"
)

func placeOrders(ordersL []orders.Order, cOK chan interface{}) {

	utils.StatusLine(3, "Processing market page")

	cDeals := make(chan deals.Deal)
	defer close(cDeals)
	go shoppinglists.ConsumeDeals(cDeals, cOK)

	for _, o := range ordersL {
		orders.Set(o)
		items.PlaceOrder(o)
		deals.ComputeDeals(o, cDeals)
	}
}

func consumePages(cPages chan []orders.Order, cOK chan interface{}) {
	for page := range cPages {
		go placeOrders(page, cOK)
		//utils.StatusLine(1, "Waiting page download")
	}
}

func getMarketPage(url string, cPages chan []orders.Order) {
	var mPage []orders.Order
	// Try until successful or hard failure handled
	for {
		res, err := utils.GetURL(url)
		if err != nil {
			fmt.Printf("Failed to get %s: %v. Skipping.\n", url, err)
			// Return empty page to satisfy progressBar/flow control
			cPages <- []orders.Order{}
			return
		}

		body, _ := io.ReadAll(res.Body)
		res.Body.Close()

		if err := json.Unmarshal(body, &mPage); err == nil && mPage != nil {
			break
		}
		// If unmarshal fails (e.g. empty body or bad json), maybe retrying helps?
		// Or just skip. For now, let's skip to avoid infinite loop on bad data.
		fmt.Printf("Failed to unmarshal %s. Skipping.\n", url)
		cPages <- []orders.Order{}
		return
	}
	cPages <- mPage
}

func producePages(lURL []string, cPages chan []orders.Order) {
	// Use worker pool to limit concurrent HTTP requests
	pool := utils.NewWorkerPool(50) // Max 50 concurrent requests

	for _, url := range lURL {
		u := url // Capture for closure
		pool.Submit(func() {
			getMarketPage(u, cPages)
		})
	}

	pool.Wait()
}

// FetchMarket will fetch the market pages from ESI
func FetchMarket() {
	items.Cleanup()
	fmt.Println("Fetching markets pages")
	cOK := make(chan interface{})
	defer close(cOK)
	cPages := make(chan []orders.Order, 100) // Buffered channel for better throughput
	defer close(cPages)

	lURL := regions.GetMarketsPagesList()
	total := len(lURL)
	go consumePages(cPages, cOK)
	go producePages(lURL, cPages)

	utils.ProgressBar(total, cOK)
}

// PrintShoppingLists is a facade method to print the n most profitable shopping lists
func PrintShoppingLists(n int) {
	shoppinglists.PrintTop(n)
}

// Terminate will clean the data structures and possibly save modifications ocurred to relevant files
func Terminate() {
	items.Cleanup()
	// shoppinglists.Cleanup() -- REMOVED for continuous updates
	// orders.Cleanup() -- REMOVED for continuous updates

	conf.Terminate()
	items.Terminate()
	regions.Terminate()
	locations.Terminate()
}
