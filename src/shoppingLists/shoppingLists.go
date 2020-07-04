package shoppinglists

import (
	"fmt"
	"os"
	"sort"
	"sync"
	"text/tabwriter"
	"time"

	"../conf"
	"../deals"
	"../locations"
	"../utils"
	"../utils/color"
)

type shopList struct {
	sellID     int64
	buyID      int64
	profit     float64
	cargoUsed  float64
	investment float64
	st         string
	deals      deals.DealsList
	selected   deals.SelectedDealsList
	itemProfit map[int]float64
}

type shopLists []*shopList

func (sl *shopList) Less(sl2 *shopList) bool {
	if sl.getProfit() == sl2.getProfit() {
		return sl.distance() < sl2.distance()
	}
	return sl.getProfit() > sl2.getProfit()
}

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

func (sl *shopList) add(d deals.Deal) {
	sl.deals = append(sl.deals, d)
}

func (sl shopList) distance() int {
	return locations.GetDistance(sl.sellID, sl.buyID)
}

func (sl shopList) key() int64 {
	return (sl.sellID * 10000000000) + sl.buyID
}

func (sl *shopList) reset() {
	sl.deals = deals.DealsList{}
	sl.selected = deals.SelectedDealsList{}
	sl.itemProfit = map[int]float64{}
	sl.profit = 0.0
	sl.cargoUsed = 0.0
}

func (sl *shopList) profitPerJump() float64 {
	return sl.getProfit() / float64(sl.distance())
}

func (sl *shopList) getProfit() float64 {
	if sl.profit > 0.0 {
		return sl.profit
	}
	cargo := conf.Cargo()
	isk := conf.MaxInvest()
	dProfit := 0.0
	qnt := 0
	sort.Sort(sl.deals)
	for _, deal := range sl.deals {
		cargo, isk, qnt, dProfit = deal.Execute(cargo, isk)
		if dProfit > 0.0 {
			sl.itemProfit[deal.GetItemID()] += dProfit
			sl.selected = append(sl.selected, deals.SelectedDeal{deal, qnt, dProfit})
			sl.profit += dProfit
		}
	}
	sl.cargoUsed = conf.Cargo() - cargo
	sl.investment = conf.MaxInvest() - isk
	for _, deal := range sl.deals {
		deal.Reset()
	}
	return sl.profit
}

func (sl shopList) String() string {
	st := ""
	less := func(i, j int) bool {
		iItemID := sl.selected[i].Deal.GetItemID()
		jItemID := sl.selected[j].Deal.GetItemID()
		return sl.itemProfit[iItemID] > sl.itemProfit[jItemID]
	}
	sort.SliceStable(sl.selected, less)
	for _, sd := range sl.selected {
		st += sd.String()
	}
	return fmt.Sprintf("\nfrom: %s", color.Fg8b(4, locations.GetName(sl.sellID))) +
		fmt.Sprintf(" to: %s", color.Fg8b(4, locations.GetName(sl.buyID))) +
		fmt.Sprintf(" %d jumps", sl.distance()) +
		st +
		fmt.Sprintf("\ntotal volume: %.2f", sl.cargoUsed) +
		fmt.Sprintf("\ninvestment: %s", color.Fg8b(1, utils.FormatCommas(sl.investment))) +
		fmt.Sprintf("\ntotal profit: %s", color.Fg8b(2, utils.FormatCommas(sl.profit))) +
		fmt.Sprintf("\nprofit per jump: %s", color.Fg8b(5, utils.FormatCommas(sl.profitPerJump())))
}

func getShopList(d deals.Deal) *shopList {
	mutex.Lock()
	defer mutex.Unlock()
	key := d.Key()
	sl, ok := listsMap[key]
	if !ok {
		sl = &shopList{d.SellLocID(), d.BuyLocID(), 0.0, 0.0, 0.0, "", deals.DealsList{}, deals.SelectedDealsList{}, map[int]float64{}}
		listsMap[key] = sl
		lists = append(lists, sl)
	}
	return sl
}

//ConsumeDeals will receive and process trade deals
func ConsumeDeals(cDeals chan deals.Deal, cOK chan bool) {
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

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 1, 2, ' ', 0)

	for i := 0; i < n; i++ {
		fmt.Fprintln(w, lists[i])
	}
	w.Flush()
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
