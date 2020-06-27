package shoppingLists

import (
    "../conf"
    "../deals"
    "../utils"
    "../locations"
    "../utils/avl"
    "../utils/color"
    "fmt"
    "sync"
)

type shopList struct {
    sellID int64
    buyID int64
    profit float64
    deals *avl.Avl
    st string
    cargoUsed float64
    investment float64
}

type dealAvlData struct {
    deal *deals.Deal
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
    //return a.sl.profitPerJump() < d.sl.profitPerJump()
    if a.sl.getProfit() == d.sl.getProfit() {
        return a.sl.distance() > d.sl.distance()
    }
    return a.sl.getProfit() < d.sl.getProfit()
}

var shopLists_m map[int64]map[int64]*shopList
var cConsumeDeals chan deals.Deal
var mutex sync.Mutex

func (s *shopList) add(d *deals.Deal) {
    ad := avl.Data(dealAvlData{d})
    s.deals.Put(&ad)
}

func (s *shopList) distance() int {
    return locations.GetDistance(s.sellID, s.buyID)
}

func (s *shopList) key() (int64, int64) {
    return s.sellID, s.buyID
}

func (s *shopList) reset() {
    s.deals = avl.NewAvl(avl.REVERSED)
    s.profit = 0.0
    s.cargoUsed = 0.0
}

func (s *shopList) profitPerJump() float64 {
    return s.getProfit()/float64(s.distance())
}

func (s *shopList) getProfit() float64 {
    if s.profit > 0.0 { return s.profit }
    it := s.deals.GetIterator()
    cargo := conf.Cargo()
    isk := conf.MaxInvest()
    sumProfit := 0.0
    strg := ""
    for it.Next() {
        adp := it.Value()
        deal := (*adp).(dealAvlData).deal
        cargo, isk, sumProfit, strg = deal.Execute(cargo, isk)
        if sumProfit > 0.0 {
            s.st += strg
        }
        s.profit += sumProfit
    }
    s.cargoUsed = conf.Cargo() - cargo
    s.investment = conf.MaxInvest() - isk
    it = s.deals.GetIterator()
    for it.Next() {
        adp := it.Value()
        deal := (*adp).(dealAvlData).deal
        deal.Reset()
    }
    return s.profit
}

func (s *shopList) String() string {
    s.st = fmt.Sprintf("\nfrom: %s", color.Fg(4, locations.Name(s.sellID)))+
        fmt.Sprintf("\tto:   %s", color.Fg(4, locations.Name(s.buyID)))+
        fmt.Sprintf("\t%d jumps", s.distance())+s.st
    s.st += fmt.Sprintf("\ntotal volume: %.2f", s.cargoUsed)
    s.st += fmt.Sprintf("\ninvestment: %s", color.Fg(1 ,utils.FormatCommas(s.investment)))
    s.st += fmt.Sprintf("\ntotal profit: %s", color.Fg(2 ,utils.FormatCommas(s.profit)))
    s.st += fmt.Sprintf("\nprofit per jump: %s", color.Fg(2 ,utils.FormatCommas(s.profitPerJump())))
    return s.st
}

func getShopList(d *deals.Deal) *shopList {
    mutex.Lock()
    defer mutex.Unlock()
    orig, dest := d.Key()
    sl, ok := shopLists_m[orig][dest]
    if !ok {
        _, ok = shopLists_m[orig]
        if !ok {
            shopLists_m[orig] = map[int64]*shopList{}
        }
        sl = &shopList{d.SellLocID(), d.BuyLocID(), 0.0, avl.NewAvl(avl.REVERSED), "", 0.0, 0.0}
        shopLists_m[orig][dest] = sl
    }
    return sl
}

func ConsumeDeals(cDeals chan *deals.Deal, cOK chan bool) {
    for d := range cDeals {
        sl := getShopList(d)
        sl.add(d)
    }
    cOK <- true
}

func PrintTop(n int) {
    fmt.Println("LISTAS")
    shopListsTree := avl.NewAvl(avl.REVERSED)
    cOK := make(chan bool)
    defer close(cOK)

    go utils.ProgressBar(len(shopLists_m), cOK)
    for _, im := range shopLists_m {
        for _, lp := range im {
            ad := avl.Data(shopListAvlData{lp})
            shopListsTree.Put(&ad)
        }
        cOK <- true
    }
    it := shopListsTree.GetIterator()
    for it.Next() {
        slp := it.Value()
        sl := (*slp).(shopListAvlData).sl

        fmt.Println(sl)

        n--
        if n == 0 { break }
    }//*/
    println()
}

func Cleanup() {
    for _, im := range shopLists_m {
        for _, lp := range im {
            lp.reset()
            lp.st = ""
        }
    }
}

func init() {
    mutex = sync.Mutex{}
    shopLists_m = map[int64]map[int64]*shopList{}
}
