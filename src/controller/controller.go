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
	"sync"

	"github.com/goccy/go-json"

	"fmt"
	"io"
)

func placeOrders(ordersL []orders.Order, cOK chan interface{}) {
	defer utils.StartTimer("Controller_PlaceOrders_Acum")()
	utils.StatusLine(3, "Processing market page")

	cDeals := make(chan deals.Deal, 100) // Buffered channel
	defer close(cDeals)
	go shoppinglists.ConsumeDealsBatch(cDeals, cOK)

	for _, o := range ordersL {
		orders.Set(o)
		items.PlaceOrder(o)
		deals.ComputeDeals(o, cDeals)
	}
}

func consumePages(cPages chan []orders.Order, cOK chan interface{}) {
	// Limit concurrency to avoid CPU thrashing based on configuration
	concurrency := conf.Threads()
	var wg sync.WaitGroup

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for page := range cPages {
				placeOrders(page, cOK)
			}
		}()
	}
	wg.Wait()
}

func getMarketPage(url string, cPages chan []orders.Order) {
	defer utils.StartTimer("Controller_GetMarketPage_Acum")()
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
	defer utils.StartTimer("Controller_FetchMarket_Total")()
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
	defer utils.StartTimer("Terminate")()
	items.Cleanup()
	// shoppinglists.Cleanup() -- REMOVED for continuous updates
	// orders.Cleanup() -- REMOVED for continuous updates

	conf.Terminate()
	items.Terminate()
	regions.Terminate()
	locations.Terminate()
}
