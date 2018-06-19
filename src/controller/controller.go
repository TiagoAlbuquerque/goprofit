package controller

import (
    "../deals"
    "../items"
    "../locations"
    "../order"
    "../regions"
    "../shoppingLists"
    "../utils"
    //    "../utils/avl"

    "encoding/json"
    "fmt"
    "io/ioutil"

    "github.com/ti/nasync"
)

func placeOrders(orders []order.Order, cOK chan bool) {

    utils.StatusIndicator("Processing market page")

    cDeals := make(chan deals.Deal)
    defer close(cDeals)
    go shoppingLists.ConsumeDeals(cDeals, cOK)

    for _, o := range orders {
        items.PlaceOrder(&o)
        deals.ComputeDeals(&o, cDeals)
    }
}

func consumePages(cPages chan []order.Order, cOK chan bool) {
    for  page := range cPages {
        placeOrders(page, cOK)
        utils.StatusIndicator("Waiting page download")
    }
}

func getMarketPages(url string, cPages chan []order.Order) {
    var mPage []order.Order
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
    cPages := make(chan []order.Order)
    defer close(cPages)
    async := nasync.New(1000, 1000)
    defer async.Close()

    lURL := regions.GetMarketsPagesList()
    total := len(lURL)
    go consumePages(cPages, cOK)

    for i := 0; i < total; i++ {
        async.Do(getMarketPages, lURL[i], cPages)
    }

    utils.ProgressBar(total, cOK)
}

func PrintShoppingLists(n int) {
    shoppingLists.PrintTop(n)
}

func Terminate() {
    items.Cleanup()
    deals.Cleanup()
    shoppingLists.Cleanup()

    items.Terminate()
    regions.Terminate()
    locations.Terminate()
}
