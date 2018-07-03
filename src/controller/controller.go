package controller

import (
    "../deals"
    "../items"
    "../locations"
    "../orders"
    "../regions"
    "../shoppingLists"
    "../utils"
    //    "../utils/avl"

    "encoding/json"
    "fmt"
    "io/ioutil"

    "github.com/ti/nasync"
)

func placeOrders(ordersL []orders.Order, cOK chan bool) {

    utils.StatusIndicator("Processing market page")

    cDeals := make(chan *deals.Deal)
    defer close(cDeals)
    go shoppingLists.ConsumeDeals(cDeals, cOK)

    for _, o := range ordersL {
        orders.Set(o)
        items.PlaceOrder(o)
        deals.ComputeDeals(o, cDeals)
    }
}

func consumePages(cPages chan []orders.Order, cOK chan bool) {
    for  page := range cPages {
        placeOrders(page, cOK)
        utils.StatusIndicator("Waiting page download")
    }
}

func getMarketPages(url string, cPages chan []orders.Order) {
    var mPage []orders.Order
    for ok := false; !ok; {
        res := utils.GetUrl(url)
        defer res.Body.Close()
        body, _ := ioutil.ReadAll(res.Body)
        json.Unmarshal(body, &mPage)
        ok = (mPage != nil)
    }
    cPages <- mPage
}

func FetchMarket() {
    items.Cleanup()
    fmt.Println("Fetching markets pages")
    cOK := make(chan bool)
    defer close(cOK)
    cPages := make(chan []orders.Order)
    defer close(cPages)
    async := nasync.New(1000, 1000)
    defer async.Close()

    lURL := regions.GetMarketsPagesList()
    total := len(lURL)
    go consumePages(cPages, cOK)

    for _, url := range lURL {
        //async.Do(getMarketPages, url, cPages)
        go getMarketPages(url, cPages)
    }

    utils.ProgressBar(total, cOK)
}

func PrintShoppingLists(n int) {
    shoppingLists.PrintTop(n)
}

func Terminate() {
    items.Cleanup()
    deals.Cleanup()
    orders.Cleanup()
    shoppingLists.Cleanup()

    items.Terminate()
    regions.Terminate()
    locations.Terminate()
}
