package shoppingLists

import (
    "../deals"
    "../utils"
    "../utils/avl"
    "fmt"
    "sync"
)

type shopList struct {
    sellID int64
    buyID int64
    deals avl.Avl
    selected avl.Avl
}

type dealAvlData struct {
    deal deals.Deal
}

func (a dealAvlData) Less (b *avl.Data) bool {
    c := (*b)
    d := c.(dealAvlData)
    return a.deal.Pm3() < d.deal.Pm3()
}

type shopListAvlData struct {
    sl *shopList
}

func (a shopListAvlData) Less (b *avl.Data) bool {
    c := (*b)
    d := c.(shopListAvlData)
    return a.sl.Profit() < d.sl.Profit()
}

var shopLists_m map[int64]map[int64]*shopList
var cConsumeDeals chan deals.Deal

func (s *shopList) add(d deals.Deal) {
    ad := avl.Data(dealAvlData{d})
    s.deals.Put(&ad)
}

func (s *shopList) key() (int64, int64) {
    return s.sellID, s.buyID
}

func (s *shopList) Profit() float64 {
    it := s.deals.GetIterator()
    for it.Next() {

    }
    out := 0.0
    return out
}

func (s *shopList) reset() {
    s.deals = avl.NewAvl(avl.REVERSED)
    s.selected = avl.NewAvl(avl.REVERSED)
}

var mutex = sync.Mutex{}
func getShopList(d deals.Deal) *shopList {
    mutex.Lock()
    defer mutex.Unlock()
    orig, dest := d.Key()
    sl, ok := shopLists_m[orig][dest]
    if !ok {
        _, ok = shopLists_m[orig]
        if !ok {
            shopLists_m[orig] = map[int64]*shopList{}
        }
        sl = &shopList{d.SellLocID(), d.BuyLocID(), avl.NewAvl(avl.REVERSED), avl.NewAvl(avl.REVERSED)}

        shopLists_m[orig][dest] = sl
    }
    return sl
}

func consumeDeals(cDeals chan deals.Deal) {
    for d := range cDeals {
        cConsumeDeals <- d
    }
}

func ConsumeDeals(cDeals chan deals.Deal) {
    for d := range cDeals {
        //sl, ok := shopLists_m[d.Key()]
        //if !ok {
        sl := getShopList(d)
            //sl = &shopList{d.SellLocID(), d.BuyLocID(), avl.NewAvl(avl.REVERSED), avl.NewAvl(avl.REVERSED)}
            //shopLists_m[sl.key()] = sl
        //}
        sl.add(d)
    }
}

func PrintTop(n int) {
    fmt.Println("LISTAS")
    shopLists_t := avl.NewAvl(avl.REVERSED)
	cOK := make(chan bool)
	go utils.ProgressBar(len(shopLists_m), cOK)
    for _, im := range shopLists_m {
        for _, lp := range im {
            ad := avl.Data(shopListAvlData{lp})
            shopLists_t.Put(&ad)
        }
        cOK <- true
    }
    it := shopLists_t.GetIterator()
    for it.Next() {

        fmt.Println(it.Value())

        n -=1
        if n == 0 { break }
    }
}

func Cleanup() {
    for _, im := range shopLists_m {
        for _, lp := range im {
            lp.reset()
        }
    }
}

func init() {
    shopLists_m = map[int64]map[int64]*shopList{}
    cConsumeDeals = make(chan deals.Deal)
    go consumeDeals(cConsumeDeals)
}
