package controller

import (
    "../conf"
    "../deals"
    "../items"
    "../utils"
    "../orders"
    "../regions"
    "../locations"
    "../shoppinglists"
    "../utils/color"
 
    //    "../utils/avl"

    "fmt"
    "io/ioutil"
    "encoding/json"
    "github.com/ti/nasync"
)

func placeOrders(ordersL []orders.Order, cOK chan bool) {

    utils.StatusIndicator(color.Fg(3, "Processing market page"))

    cDeals := make(chan *deals.Deal)
    defer close(cDeals)
    go shoppinglists.ConsumeDeals(cDeals, cOK)

    for _, o := range ordersL {
        orders.Set(o)
        items.PlaceOrder(o)
        deals.ComputeDeals(o, cDeals)
    }
}

func consumePages(cPages chan []orders.Order, cOK chan bool) {
    for  page := range cPages {
        placeOrders(page, cOK)
        utils.StatusIndicator(color.Fg(1, "Waiting page download"))
    }
}

func getMarketPage(url string, cPages chan []orders.Order) {
    var mPage []orders.Order
    for ok := false; !ok; {
        res := utils.GetURL(url)
        defer res.Body.Close()
        body, _ := ioutil.ReadAll(res.Body)
        json.Unmarshal(body, &mPage)
        ok = (mPage != nil)
    }
    cPages <- mPage
}


func consumeMarketPages(cURL chan string, cPages chan []orders.Order) {
    for url := range cURL {
        getMarketPage(url, cPages)
    }

}

func getMarketPages(lURL []string, cPages chan []orders.Order) {
    cURL := make(chan string)
    defer close(cURL)
    for i := 0; i < 210; i++ {
        go consumeMarketPages(cURL, cPages)
    }
    for _, url := range lURL {
        cURL <- url
    }
}
//FetchMarket will fetch the market pages from ESI 
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

    //go getMarketPages(lURL, cPages)

    for _, url := range lURL {
        async.Do(getMarketPage, url, cPages)
      //  go getMarketPage(url, cPages)
    }

    utils.ProgressBar(total, cOK)
}

//PrintShoppingLists is a facade method to print the n most profitable shopping lists
func PrintShoppingLists(n int) {
    shoppinglists.PrintTop(n)
}

//Terminate will clean the data structures and possibly save modifications ocurred to relevant files
func Terminate() {
    items.Cleanup()
    deals.Cleanup()
    orders.Cleanup()
    shoppinglists.Cleanup()

    conf.Terminate()
    items.Terminate()
    regions.Terminate()
    locations.Terminate()

}
