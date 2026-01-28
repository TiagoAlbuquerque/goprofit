package shoppinglists

import (
	"fmt"
	"os"
	"sort"
	"sync"
	"text/tabwriter"
	"time"

	"goprofit/conf"
	"goprofit/deals"
	"goprofit/items"
	"goprofit/locations"
	"goprofit/orders"
	"goprofit/utils"
	"goprofit/utils/color"
)

type shopList struct {
	sellID     int64
	buyID      int64
	profit     float64
	cargoUsed  float64
	investment float64
	deals      deals.List
	selected   deals.SelectedList
	itemProfit map[int]float64
	dealKeys   map[string]int // Map of "buyOrderID-sellOrderID" -> index in deals slice for O(1) deduplication
}

type shopLists []*shopList

func (sl *shopList) Less(sl2 *shopList) bool {
	//if sl.getProfit() == sl2.getProfit() {
	//	return sl.distance() < sl2.distance()
	//}
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
	defer utils.StartTimer("ShoppingList_Add_Acum")()
	// O(1) deduplication using map
	key := fmt.Sprintf("%d-%d", d.GetBuyOrderID(), d.GetSellOrderID())
	if idx, exists := sl.dealKeys[key]; exists {
		// Update existing deal round
		sl.deals[idx].Round = currentRound
		return
	}
	d.Round = currentRound
	sl.dealKeys[key] = len(sl.deals)
	sl.deals = append(sl.deals, d)
}

func (sl shopList) distance() int {
	return locations.GetDistance(sl.sellID, sl.buyID)
}

func (sl shopList) key() int64 {
	return (sl.sellID * 10000000000) + sl.buyID
}

func (sl *shopList) reset() {
	sl.deals = deals.List{}
	sl.selected = deals.SelectedList{}
	sl.itemProfit = map[int]float64{}
	sl.profit = 0.0
	sl.cargoUsed = 0.0
}

func (sl *shopList) profitPerJump() float64 {
	return sl.getProfit() / float64(sl.distance())
}

func (sl *shopList) selectDeal(deal deals.Deal, res *deals.Resources) {
	isSelected, sDeal := deal.Execute(res)
	if isSelected {
		sl.itemProfit[deal.GetItemID()] += sDeal.Profit
		sl.selected = append(sl.selected, sDeal)
		sl.profit += sDeal.Profit
	}
}

func (sl *shopList) getProfit() float64 {
	defer utils.StartTimer("ShoppingList_GetProfit_Acum")()
	if sl.profit > 0.0 {
		return sl.profit
	}
	res := deals.Resources{Cargo: conf.Cargo(), Isk: conf.MaxInvest()}
	sort.Sort(sl.deals)
	for _, deal := range sl.deals {
		sl.selectDeal(deal, &res)
	}
	sl.cargoUsed = conf.Cargo() - res.Cargo
	sl.investment = conf.MaxInvest() - res.Isk
	for _, deal := range sl.deals {
		deal.Reset()
	}
	return sl.profit
}
func (sl shopList) wappString() string {
	return fmt.Sprintf("\nfrom: %s", locations.GetName(sl.sellID)) +
		fmt.Sprintf(" to: %s", locations.GetName(sl.buyID)) +
		fmt.Sprintf(" %d jumps", sl.distance()) +
		fmt.Sprintf("\ntotal volume: %.2f", sl.cargoUsed) +
		fmt.Sprintf("\ninvestment: %s", utils.FormatCommas(sl.investment)) +
		fmt.Sprintf("\ntotal profit: %s", utils.FormatCommas(sl.profit)) +
		fmt.Sprintf("\nprofit per jump: %s", utils.FormatCommas(sl.profitPerJump()))
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

func getShopList(key int64, sellLoc int64, buyLoc int64) *shopList {
	sl, ok := listsMap[key]
	if !ok {
		sl = &shopList{
			sellID:     sellLoc,
			buyID:      buyLoc,
			profit:     0.0,
			cargoUsed:  0.0,
			investment: 0.0,
			deals:      deals.List{},
			selected:   deals.SelectedList{},
			itemProfit: map[int]float64{},
			dealKeys:   map[string]int{},
		}
		listsMap[key] = sl
		lists = append(lists, sl)
	}
	return sl
}

// ConsumeDeals will receive and process trade deals
func ConsumeDeals(cDeals chan deals.Deal, cOK chan interface{}) {
	defer utils.StartTimer("ShoppingList_ConsumeDeals_Total")()
	for d := range cDeals {
		mutex.Lock()
		sl := getShopList(d.Key(), d.SellLocID(), d.BuyLocID())
		sl.add(d)
		mutex.Unlock()
	}
	cOK <- true
}

type dealCtx struct {
	d       deals.Deal
	key     int64
	sellLoc int64
	buyLoc  int64
}

type ingestRequest struct {
	batch []dealCtx
	cOK   chan interface{}
}

// Buffer upped to avoid any backpressure from fetchers
var ingestChan = make(chan ingestRequest, 5000)

func aggregator() {
	for req := range ingestChan {
		mutex.Lock()
		for _, ctx := range req.batch {
			sl := getShopList(ctx.key, ctx.sellLoc, ctx.buyLoc)
			sl.add(ctx.d)
		}
		mutex.Unlock()

		// Signal completion to the worker
		req.cOK <- true
	}
}

// ConsumeDealsBatch will collect all deals from the channel and send to aggregator
func ConsumeDealsBatch(cDeals chan deals.Deal, cOK chan interface{}) {
	defer utils.StartTimer("ShoppingList_ConsumeDealsBatch_Total")()

	var batch []dealCtx

	// 1. Collect all deals and PRE-CALCULATE keys/locs
	for d := range cDeals {
		batch = append(batch, dealCtx{
			d:       d,
			key:     d.Key(),
			sellLoc: d.SellLocID(),
			buyLoc:  d.BuyLocID(),
		})
	}

	// 2. Send to Aggregator
	if len(batch) > 0 {
		ingestChan <- ingestRequest{batch: batch, cOK: cOK}
	} else {
		cOK <- true
	}
}

// PrintTop will print the top n most profitable shopping lists
func PrintTop(n int) {
	fmt.Println("LISTAS")

	start := time.Now()
	utils.Top(lists)
	utils.StatusLine(15, "sorted in: "+fmt.Sprint(time.Now().Sub(start)))

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 1, 2, ' ', 0)

	count := n
	if count > len(lists) {
		count = len(lists)
	}

	for i := count - 1; i >= 0; i-- {
		fmt.Fprintln(w, lists[i])
	}
	w.Flush()
	println()
}

var lastRoundDTOs []ShoppingListDTO

// Cleanup will reset the shopping lists computed on the last round
// For continuous updates, we might NOT want to clear lists, but maybe reset some state?
// The plan says "Cleanup() to be empty or removed".
func Cleanup() {
	// No-op for continuous persistence
}

var currentRound int = 0

func NextRound() {
	mutex.Lock()
	defer mutex.Unlock()
	currentRound++
}

func Prune() {
	defer utils.StartTimer("ShoppingList_Prune")()
	mutex.Lock()
	defer mutex.Unlock()

	// Remove deals older than 1 round (currentRound - deal.Round > 1)
	for _, sl := range lists {
		var activeDeals deals.List
		newDealKeys := make(map[string]int)
		for _, d := range sl.deals {
			if currentRound-d.Round <= 1 {
				key := fmt.Sprintf("%d-%d", d.GetBuyOrderID(), d.GetSellOrderID())
				newDealKeys[key] = len(activeDeals)
				activeDeals = append(activeDeals, d)
			}
		}
		sl.deals = activeDeals
		sl.dealKeys = newDealKeys

		// Reset stats since deals changed
		sl.profit = 0
		sl.cargoUsed = 0
		sl.itemProfit = map[int]float64{}
		sl.selected = deals.SelectedList{}
	}
}

type ShoppingListDTO struct {
	From       string        `json:"From"`
	To         string        `json:"To"`
	Jumps      int           `json:"Jumps"`
	Profit     string        `json:"Profit"`
	Investment string        `json:"Investment"`
	Volume     string        `json:"Volume"`
	ROI        string        `json:"ROI"`
	Items      []ListItemDTO `json:"Items"`
}

type ListItemDTO struct {
	Name      string `json:"Name"`
	Quantity  int    `json:"Quantity"`
	BuyPrice  string `json:"BuyPrice"`
	SellPrice string `json:"SellPrice"`
	Profit    string `json:"Profit"`
}

func getTopDTOInternal(n int) []ShoppingListDTO {
	defer utils.StartTimer("ShoppingList_GetTopDTO")()
	utils.Top(lists)

	count := n
	if count > len(lists) {
		count = len(lists)
	}

	var result []ShoppingListDTO
	for i := 0; i < count; i++ {
		sl := lists[i]
		profit := sl.getProfit()
		investment := sl.investment
		if investment == 0 {
			investment = 1
		}

		var listItems []ListItemDTO
		for _, sd := range sl.selected {
			itmID := sd.Deal.GetItemID()
			itmName := items.Get(itmID).Name

			bFor := orders.Get(sd.Deal.GetSellOrderID()).Price
			sFor := orders.Get(sd.Deal.GetBuyOrderID()).Price

			itemDTO := ListItemDTO{
				Name:      itmName,
				Quantity:  sd.Qnt,
				BuyPrice:  utils.KMB(bFor),
				SellPrice: utils.KMB(sFor),
				Profit:    utils.KMB(sd.Profit),
			}
			listItems = append(listItems, itemDTO)
		}

		dto := ShoppingListDTO{
			From:       locations.GetName(sl.sellID),
			To:         locations.GetName(sl.buyID),
			Jumps:      sl.distance(),
			Profit:     utils.FormatCommas(profit),
			Investment: utils.FormatCommas(sl.investment),
			Volume:     fmt.Sprintf("%.2f", sl.cargoUsed),
			ROI:        fmt.Sprintf("%.0f", (profit/investment)*100),
			Items:      listItems,
		}
		result = append(result, dto)
	}
	return result
}

func GetTopDTO(n int) []ShoppingListDTO {
	mutex.Lock()
	defer mutex.Unlock()
	return getTopDTOInternal(n)
}

func init() {
	//cSorting = make(chan bool)
	mutex = sync.Mutex{}
	listsMap = map[int64]*shopList{}
	lists = shopLists{}

	go aggregator()
}
