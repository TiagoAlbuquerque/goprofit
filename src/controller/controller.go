package controller

import (
    "../items"
    _ "../deals"
    "../order"
    "../regions"
    "../utils"
//    "../utils/avl"

    "fmt"
    "encoding/json"
    "io/ioutil"
    "github.com/ti/nasync"
)


func placeOrders(orders []order.Order) {
    for _, o := range orders {
        items.PlaceOrder(&o)
//        deals.ComputeDeals(&o)
    }
}

func consumePages(cPages chan []order.Order, cOK chan bool, total int) {
    items.Cleanup()
    for i := 0; i < total; i++ {
        placeOrders(<-cPages)
        cOK <- true
    }
}

func getMarketPages(url string, cPages chan []order.Order){
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

func FetchMarket(){
    fmt.Println("Fetching markets pages")
    lURL := regions.GetMarketsPagesList()
    cOK := make(chan bool)
    cPages := make(chan []order.Order)
    total := len(lURL)
    go consumePages(cPages, cOK, total)
    async := nasync.New(1000, 1000)
    defer async.Close()
    for i := 0; i < total; i++ {
        async.Do(getMarketPages, lURL[i], cPages)
    }
    utils.ProgressBar(total, cOK)
}
func Terminate() {
    items.Terminate()
    regions.Terminate()
}