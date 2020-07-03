package deals

import (
	"../conf"
	"../items"
	"../orders"
	"../utils"
	"../utils/color"

	"fmt"
	"math"
)

//Deal is astructure to couple a sell and a buy order of a specific item
type Deal struct {
	itemID                  int
	buyOrderID, sellOrderID int64
}

type DealsList []Deal

func (dl DealsList) Len() int {
	return len(dl)
}
func (dl DealsList) Less(i, j int) bool {
	return dl[i].Pm3() > dl[j].Pm3()
}
func (dl DealsList) Swap(i, j int) {
	dl[i], dl[j] = dl[j], dl[i]
}

var deals DealsList

//GetItemID will return the item ID of this deal
func (d Deal) GetItemID() int {
	return d.itemID
}

//Key will produce key for the marketlist at witch this deal is to be inserted
func (d Deal) Key() int64 {
	return (orders.Get(d.sellOrderID).SystemID * 10000000000) + orders.Get(d.buyOrderID).SystemID
}

//SellLocID will return the sell location ID for the current deal
func (d Deal) SellLocID() int64 {
	return int64(orders.Get(d.sellOrderID).SystemID)
}

//BuyLocID will return the buy location ID for the current deal
func (d Deal) BuyLocID() int64 {
	return int64(orders.Get(d.buyOrderID).SystemID)
}

//Pm3 will return the profit amount normalized by cubic meter ocupied by the item
func (d Deal) Pm3() float64 {
	itm := items.Get(d.itemID)
	prf := d.profitPerUnit()
	vol := float64(itm.Volume)
	out := prf / vol

	return out
}

func min(nums ...int) int {
	out := math.Inf(1)
	for _, val := range nums {
		out = math.Min(float64(out), float64(val))
	}
	return int(out)
}

//amount that is available compose in buy/sell order
func (d Deal) amountAvailable() int {
	bo := orders.Get(d.buyOrderID)
	so := orders.Get(d.sellOrderID)

	out := min(bo.OrderRemain(), so.OrderRemain())

	return out
}

//amount that can be bought
func (d Deal) amountIsk(iskAvail float64) int {
	sop := orders.Get(d.sellOrderID).Price // sell order price
	return int(math.Floor(iskAvail / sop))
}

//amount that can fit in cargo
func (d Deal) amountCargo(cargo float64) int {
	itmVol := items.Get(d.itemID).Volume
	return int(math.Floor(cargo / itmVol))
}

func tax() float64 {
	return 1.0 - conf.Tax()
}

func (d Deal) profitPerUnit() float64 {
	bo := orders.Get(d.buyOrderID)
	so := orders.Get(d.sellOrderID)
	ppu := (bo.Price * tax()) - so.Price
	return ppu
}

func (d Deal) profitQnt(qnt int) float64 {
	ppu := d.profitPerUnit()
	out := float64(qnt) * ppu

	return out
}

func (d Deal) getQuantity(cargo, isk float64) int {
	return min(d.amountAvailable(), d.amountCargo(cargo), d.amountIsk(isk))
}

//Execute will execute the deal for as many item as its availabe to trade in this deal, can be stored in the ships cargobay, and have enough isk to purchase
func (d *Deal) Execute(cargo, isk float64) (float64, float64, float64, string) {
	itm := items.Get(d.itemID)
	bo := orders.Get(d.buyOrderID)
	so := orders.Get(d.sellOrderID)

	itmVol := itm.Volume
	itmName := itm.Name

	qnt := d.getQuantity(cargo, isk)
	bo.Execute(qnt)
	so.Execute(qnt)

	vol := float64(qnt) * itmVol
	cargo -= vol
	cost := float64(qnt) * so.Price
	isk -= cost

	bFor := so.Price
	sFor := bo.Price
	profit := d.profitQnt(qnt)

	strg := fmt.Sprintf("\n%d \t%s \tb: %s \ts: %s \tp: %s",
		qnt,
		itmName,
		color.Fg8b(3, utils.KMB(bFor)),
		color.Fg8b(6, utils.KMB(sFor)),
		color.Fg8b(2, utils.KMB(profit)))

	return cargo, isk, profit, strg
}

//Reset will restore the deal to unexecuted state
func (d *Deal) Reset() {
	bo := orders.Get(d.buyOrderID)
	bo.Reset()
	//orders.Set(bo)

	so := orders.Get(d.sellOrderID)
	so.Reset()
	//orders.Set(so)
}

func (d Deal) valid() bool {
	bo := orders.Get(d.buyOrderID)

	if d.profitPerUnit() <= 0.0 ||
		d.Pm3() < conf.Minpm3() ||
		bo.MinVolume > 1 ||
		items.Get(d.itemID).IsOfficer() {
		return false
	}
	return true
}

func makeDeal(itmID int, boID int64, soID int64, cDeals chan Deal) {
	d := Deal{itmID, boID, soID}
	if d.valid() {
		deals = append(deals, d)
		cDeals <- d
	}
}

func computeBuyOrder(bOrder orders.Order, cDeals chan Deal) {
	itm := items.Get(bOrder.ItemID)

	for _, sOrderID := range itm.SellOrders {
		makeDeal(itm.ItemID, bOrder.OrderID, sOrderID, cDeals)
	}
}

func computeSellOrder(sOrder orders.Order, cDeals chan Deal) {
	itm := items.Get(sOrder.ItemID)

	for _, bOrderID := range itm.BuyOrders {
		makeDeal(itm.ItemID, bOrderID, sOrder.OrderID, cDeals)
	}
}

//Cleanup will discard all deals stored
func Cleanup() {
	deals = []Deal{}
}

//ComputeDeals will produce Deals based on received orders
func ComputeDeals(o orders.Order, cDeals chan Deal) {
	if o.IsBuyOrder {
		computeBuyOrder(o, cDeals)
	} else {
		computeSellOrder(o, cDeals)
	}
}
