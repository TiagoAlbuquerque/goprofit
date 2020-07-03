package shoppinglists

import (
	"fmt"
	"sync"
	"time"

	"../conf"
	"../deals"
	"../locations"
	"../utils"
	"../utils/avl"
	"../utils/color"
)

type shopList struct {
	sellID     int64
	buyID      int64
	profit     float64
	deals      *avl.Avl
	st         string
	cargoUsed  float64
	investment float64
}

func (sl *shopList) Less(sl2 *shopList) bool {
	if sl.getProfit() == sl2.getProfit() {
		return sl.distance() < sl2.distance()
	}
	return sl.getProfit() > sl2.getProfit()
}

type dealAvlData struct {
	deal *deals.Deal
}

func (a dealAvlData) Less(b *avl.Data) bool {
	c := (*b)
	d := c.(dealAvlData)
	return a.deal.Pm3() < d.deal.Pm3()
}

type shopLists []*shopList

func (sls shopLists) Len() int {
	return len(sls)
}
func (sls shopLists) Less(a, b int) bool {
	return sls[a].Less(sls[b])
}
func (sls shopLists) Swap(a, b int) {
	sls[a], sls[b] = sls[b], sls[a]
}

var listsMap map[int64]*shopList
var lists shopLists
var cConsumeDeals chan deals.Deal
var mutex sync.Mutex
var cSorting chan bool

func (sl *shopList) add(d *deals.Deal) {
	ad := avl.Data(dealAvlData{d})
	sl.deals.Put(&ad)
}

func (sl shopList) distance() int {
	return locations.GetDistance(sl.sellID, sl.buyID)
}

func (sl shopList) key() int64 {
	return (sl.sellID * 10000000000) + sl.buyID
}

func (sl *shopList) reset() {
	sl.deals = avl.NewAvl(avl.REVERSED)
	sl.profit = 0.0
	sl.cargoUsed = 0.0
	sl.st = ""
}

func (sl *shopList) profitPerJump() float64 {
	return sl.getProfit() / float64(sl.distance())
}

func (sl *shopList) getProfit() float64 {
	if sl.profit > 0.0 {
		return sl.profit
	}
	it := sl.deals.GetIterator()
	cargo := conf.Cargo()
	isk := conf.MaxInvest()
	sumProfit := 0.0
	strg := ""
	for it.Next() {
		adp := it.Value()
		deal := (*adp).(dealAvlData).deal
		cargo, isk, sumProfit, strg = deal.Execute(cargo, isk)
		if sumProfit > 0.0 {
			sl.st += strg
		}
		sl.profit += sumProfit
	}
	sl.cargoUsed = conf.Cargo() - cargo
	sl.investment = conf.MaxInvest() - isk
	it = sl.deals.GetIterator()
	for it.Next() {
		adp := it.Value()
		deal := (*adp).(dealAvlData).deal
		deal.Reset()
	}
	//cSorting <- true
	return sl.profit
}

func (sl shopList) String() string {
	return fmt.Sprintf("\nfrom: %s", color.Fg8b(4, locations.GetName(sl.sellID))) +
		fmt.Sprintf("\tto:   %s", color.Fg8b(4, locations.GetName(sl.buyID))) +
		fmt.Sprintf("\t%d jumps", sl.distance()) +
		sl.st +
		fmt.Sprintf("\ntotal volume: %.2f", sl.cargoUsed) +
		fmt.Sprintf("\ninvestment: %s", color.Fg8b(1, utils.FormatCommas(sl.investment))) +
		fmt.Sprintf("\ntotal profit: %s", color.Fg8b(2, utils.FormatCommas(sl.profit))) +
		fmt.Sprintf("\nprofit per jump: %s", color.Fg8b(5, utils.FormatCommas(sl.profitPerJump())))
}

func getShopList(d *deals.Deal) *shopList {
	mutex.Lock()
	defer mutex.Unlock()
	key := d.Key()
	sl, ok := listsMap[key]
	if !ok {
		sl = &shopList{d.SellLocID(), d.BuyLocID(), 0.0, avl.NewAvl(avl.REVERSED), "", 0.0, 0.0}
		listsMap[key] = sl
		lists = append(lists, sl)
	}
	return sl
}

//ConsumeDeals will receive and process trade deals
func ConsumeDeals(cDeals chan *deals.Deal, cOK chan bool) {
	for d := range cDeals {
		sl := getShopList(d)
		sl.add(d)
	}
	cOK <- true
}

//PrintTop will print the top n most profitable shopping lists
func PrintTop(n int) {
	fmt.Println("LISTAS")

	start := time.Now()
	utils.Top(lists)
	utils.StatusLine("sorted in: " + fmt.Sprint(time.Now().Sub(start)))

	for i := 0; i < n; i++ {
		fmt.Println(lists[i])
	}
	println()
}

//Cleanup will reset the shopping lists computed on the last round
func Cleanup() {
	for _, sl := range lists {
		sl.reset()
	}
}
func init() {
	//cSorting = make(chan bool)
	mutex = sync.Mutex{}
	listsMap = map[int64]*shopList{}
	lists = shopLists{}
}
