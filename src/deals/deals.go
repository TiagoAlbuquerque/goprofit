package deals

import (
    "../items"
    "../order"
    "../utils/avl"
    _ "fmt"
    "math"
)

type Deal struct{
    item *items.Item
    buyOrder, sellOrder *order.Order
}

var deals []Deal

func (d *Deal) Key() (int64, int64) {
    return d.sellOrder.LocationID, d.buyOrder.LocationID
}

func (d *Deal) SellLocID() int64{
    return d.sellOrder.LocationID
}

func (d *Deal) BuyLocID() int64{
    return d.buyOrder.LocationID
}

func (d *Deal) Pm3() float64 {
    prf := d.profitPerUnit()
    vol := float64(d.item.Volume)
    out := prf/vol

    return out
}

func min(a, b int) int {
    if a < b { return a }
    return b
}

func (d *Deal) amount() int {
    out := min(d.buyOrder.OrderRemain(), d.sellOrder.OrderRemain())

    if out < d.buyOrder.MinVolume { out = 0 }
    if out < d.sellOrder.MinVolume { out = 0 }

    return out
}

func (d *Deal) amountCargo(cargo float64) int {
    out := d.amount()
    out = min(out, int(math.Floor(cargo/d.item.Volume)))
    return out
}

func tax() float64 {
    return 1-0.01
}

func (d *Deal) profitPerUnit() float64 {
    ppu := (d.buyOrder.Price*tax()) - d.sellOrder.Price
    return ppu
}

func (d *Deal) profitQnt(qnt int) float64 {
    ppu := d.profitPerUnit()
    out := float64(qnt)*ppu

    return out
}

func (d *Deal) Execute(cargo float64) (float64, float64) {
    itmVol := d.item.Volume
    qnt := d.amountCargo(cargo)

    d.buyOrder.Execute(qnt)
    d.sellOrder.Execute(qnt)

    vol := float64(qnt)*itmVol
    cargo -= vol

    profit := d.profitQnt(qnt)

    return cargo, profit
}

func makeDeal(item *items.Item, bOrder *order.Order, sOrder *order.Order, cDeals chan Deal) bool {
    d := Deal{item, bOrder, sOrder}
    if d.profitPerUnit() > 0 {
        deals = append(deals, d)
        cDeals <- d
        return true
    }
    return false
}

func computeBuyOrder(item *items.Item, bOrder *order.Order, sAvl *avl.Avl, cDeals chan Deal) {
    iter := sAvl.GetIterator()

    for iter.Next() {
        sOrder := (*iter.Value()).(items.OrderAvlData).Order
        if !makeDeal(item, bOrder, sOrder, cDeals) {
            break
        }
    }
}

func computeSellOrder(item *items.Item, sOrder *order.Order, bAvl *avl.Avl, cDeals chan Deal) {
    iter := bAvl.GetIterator()

    for iter.Next() {
        bOrder := (*iter.Value()).(items.OrderAvlData).Order
        if !makeDeal(item, bOrder, sOrder, cDeals) {
            break
        }
    }
}

func Cleanup() {
    deals = []Deal{}
}

func ComputeDeals(o *order.Order, cDeals chan Deal) {
    itID := o.ItemID
    item := items.GetItem(itID)
    if o.IsBuyOrder {
        computeBuyOrder(item, o, item.Sell_orders, cDeals)
    } else {
        computeSellOrder(item, o, item.Buy_orders, cDeals)
    }

}
