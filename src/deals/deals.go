package deals

import (
	"../conf"
	"../items"
	"../orders"

	"fmt"
	"math"
)

type Deal struct {
	item                int
	buyOrder, sellOrder int64
}

var deals []Deal

func (d *Deal) Key() int64 {
	return (orders.Get(d.sellOrder).SystemID * 10000000000) + orders.Get(d.buyOrder).SystemID
}

func (d *Deal) SellLocID() int64 {
	return int64(orders.Get(d.sellOrder).SystemID)
}

func (d *Deal) BuyLocID() int64 {
	return int64(orders.Get(d.buyOrder).SystemID)
}

func (d *Deal) Pm3() float64 {
	itm := items.Get(d.item)
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
func (d *Deal) amountAvailable() int {
	bo := orders.Get(d.buyOrder)
	so := orders.Get(d.sellOrder)

	out := min(bo.OrderRemain(), so.OrderRemain())

	return out
}

//amount that can be bought
func (d *Deal) amountIsk(iskAvail float64) int {
	sop := orders.Get(d.sellOrder).Price // sell order price
	return int(math.Floor(iskAvail / sop))
}

//amount that can fit in cargo
func (d *Deal) amountCargo(cargo float64) int {
	itmVol := items.Get(d.item).Volume
	return int(math.Floor(cargo / itmVol))
}

func tax() float64 {
	return 1.0 - conf.Tax()
}

func (d *Deal) profitPerUnit() float64 {
	bo := orders.Get(d.buyOrder)
	so := orders.Get(d.sellOrder)
	ppu := (bo.Price * tax()) - so.Price
	return ppu
}

func (d *Deal) profitQnt(qnt int) float64 {
	ppu := d.profitPerUnit()
	out := float64(qnt) * ppu

	return out
}

func (d *Deal) getQuantity(cargo, isk float64) int {
	return min(d.amountAvailable(), d.amountCargo(cargo), d.amountIsk(isk))
}

func (d *Deal) Execute(cargo, isk float64) (float64, float64, float64, string) {
	itm := items.Get(d.item)
	bo := orders.Get(d.buyOrder)
	so := orders.Get(d.sellOrder)

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

	strg := fmt.Sprintf("\n%d\t%s \tbuy for: %.2f \tsell for: %.2f \tprofit: %.2f",
		qnt,
		itmName,
		bFor,
		sFor,
		profit)

	return cargo, isk, profit, strg
}

func (d *Deal) Reset() {
	bo := orders.Get(d.buyOrder)
	bo.Reset()
	//orders.Set(bo)

	so := orders.Get(d.sellOrder)
	so.Reset()
	//orders.Set(so)
}

func (d *Deal) valid() bool {
	bo := orders.Get(d.buyOrder)

	if d.profitPerUnit() <= 0.0 ||
		d.Pm3() < conf.Minpm3() ||
		bo.MinVolume > 1 ||
		items.Get(d.item).IsOfficer() {
		return false
	}
	return true
}

func makeDeal(itmID int, boID int64, soID int64, cDeals chan *Deal) {
	d := Deal{itmID, boID, soID}
	if d.valid() {
		deals = append(deals, d)
		cDeals <- &d
	}
}

func computeBuyOrder(bOrder orders.Order, cDeals chan *Deal) {
	itm := items.Get(bOrder.ItemID)

	for _, sOrder := range itm.SellOrders {
		makeDeal(itm.ItemID, bOrder.OrderID, sOrder, cDeals)
	}
}

func computeSellOrder(sOrder orders.Order, cDeals chan *Deal) {
	itm := items.Get(sOrder.ItemID)

	for _, bOrder := range itm.BuyOrders {
		makeDeal(itm.ItemID, bOrder, sOrder.OrderID, cDeals)
	}
}

func Cleanup() {
	deals = []Deal{}
}

func ComputeDeals(o orders.Order, cDeals chan *Deal) {
	if o.IsBuyOrder {
		computeBuyOrder(o, cDeals)
	} else {
		computeSellOrder(o, cDeals)
	}
}
