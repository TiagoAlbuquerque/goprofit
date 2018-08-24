package deals

import (
    "../items"
    "../orders"
//    "../utils/avl"
    "fmt"
    "math"
)

type Deal struct{
    item int
    buyOrder, sellOrder int64
}

var deals []Deal

func (d *Deal) Key() (int64, int64) {
    return orders.Get(d.sellOrder).LocationID, orders.Get(d.buyOrder).LocationID
}

func (d *Deal) SellLocID() int64{
    return orders.Get(d.sellOrder).LocationID
}

func (d *Deal) BuyLocID() int64{
    return orders.Get(d.buyOrder).LocationID
}

func (d *Deal) Pm3() float64 {
    itm := items.Get(d.item)
    prf := d.profitPerUnit()
    vol := float64(itm.Volume)
    out := prf/vol

    return out
}

func min(a, b int) int {
    if a < b { return a }
    return b
}

func (d *Deal) amount() int {
    bo := orders.Get(d.buyOrder)
    so := orders.Get(d.sellOrder)

    out := min(bo.OrderRemain(), so.OrderRemain())

    if bo.MinVolume > 1 { out = 0 }

    return out
}

func (d *Deal) amountCargo(cargo float64) int {
    out := d.amount()
    itm := items.Get(d.item)
    out = min(out, int(math.Floor(cargo/itm.Volume)))
    return out
}

func tax() float64 {
    return 1-0.01
}

func (d *Deal) profitPerUnit() float64 {
    bo := orders.Get(d.buyOrder)
    so := orders.Get(d.sellOrder)
    ppu := (bo.Price*tax()) - so.Price
    return ppu
}

func (d *Deal) profitQnt(qnt int) float64 {
    ppu := d.profitPerUnit()
    out := float64(qnt)*ppu

    return out
}

func (d *Deal) Execute(cargo float64) (float64, float64, string) {
    itm := items.Get(d.item)
    bo := orders.Get(d.buyOrder)
    so := orders.Get(d.sellOrder)

    itmVol := itm.Volume
    itmName := itm.Name

    qnt := d.amountCargo(cargo)

    bo.Execute(qnt)
    so.Execute(qnt)

    vol := float64(qnt)*itmVol
    cargo -= vol

    bFor := so.Price
    sFor := bo.Price
    profit := d.profitQnt(qnt)

    strg := fmt.Sprintf("\n%d\t%s \tbuy for: %.2f \tsell for: %.2f \tprofit: %.2f",
                        qnt,
                        itmName,
                        bFor,
                        sFor,
                        profit)

    return cargo, profit, strg
}

func (d *Deal) Reset() {
    bo := orders.Get(d.buyOrder)
    bo.Reset()
    orders.Set(bo)

    so := orders.Get(d.sellOrder)
    so.Reset()
    orders.Set(so)
}

func makeDeal(itmID int, boID int64, soID int64, cDeals chan *Deal) bool {
    d := Deal{itmID, boID, soID}
    if d.profitPerUnit() > 0.0 && d.Pm3() > 100000{
        deals = append(deals, d)
        cDeals <- &d
        return true
    }
    return false
}

func computeBuyOrder(bOrder orders.Order, cDeals chan *Deal) {
    itm := items.Get(bOrder.ItemID)

    for _, sOrder := range itm.Sell_orders {
        makeDeal(itm.ItemID, bOrder.OrderID, sOrder, cDeals)
    }
}

func computeSellOrder(sOrder orders.Order, cDeals chan *Deal) {
    itm := items.Get(sOrder.ItemID)

    for _, bOrder := range itm.Buy_orders {
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
