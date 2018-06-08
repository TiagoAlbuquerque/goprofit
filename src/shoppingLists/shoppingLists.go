package shoppingLists

import (
    "../deals"
    "../utils/avl"
    "fmt"
)

type shopList struct {
    sellID int64
    buyID int64
    deals_l avl.Avl
    selected_l avl.Avl
}

var shopLists_m map[string]*shopList
var shopLists_t avl.Avl
var cConsumeDeals chan deals.Deal

func (s *shopList) add(d deals.Deal) {
    s.deals_l.Put(d)
}

func (s *shopList) Key() string {
    return fmt.Sprintf("%d >> %d", s.sellID, s.buyID)
}

func ConsumeDeals(cDeals chan deals.Deal) {
    for d := range cDeals {
        cConsumeDeals <- d
    }
}
func consumeDeals(cDeals chan deals.Deal) {
    for d := range cDeals {
        sl, ok := shopLists_m[d.Key()]
        if !ok {
            sl = &shopList{d.SellLocID(), d.BuyLocID(), avl.NewAvl(avl.REVERSED), avl.NewAvl(avl.REVERSED)}
            shopLists_m[sl.Key()] = sl
        }
        sl.add(d)
    }
}
func PrintTop(n int) {
    fmt.Println("LISTAS")
}
func Cleanup() {
    for _, lp := range shopLists_m {
        lp.deals_l = []deals.Deal{}
        lp.selected_l = []deals.Deal{}
    }
}
func init() {
    shopLists_m = map[string]*shopList{}
    shopLists_t = avl.NewAvl(avl.REVERSED)
    cConsumeDeals = make(chan deals.Deal)
    go consumeDeals(cConsumeDeals)
}
