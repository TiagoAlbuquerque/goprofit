package deals

import (
    _ "fmt"
    "../items"
    "../order"
    "../utils/avl"
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
    prf := d.Profit()
    amt := d.amount()
    vol := float64(amt)*float64(d.item.Volume)
    out := prf/vol

    return out
}

func min(a, b int) int {
    if a < b { return a }
    return b
}

func (d *Deal) amount() int {
    out := min(d.buyOrder.VolumeRemain, d.sellOrder.VolumeRemain)

    if out < d.buyOrder.MinVolume { out = 0 }
    if out < d.sellOrder.MinVolume { out = 0 }

    return out
}

func (d *Deal) Profit() float64 {
    tax := 1-0.01
    amt := d.amount()
    ppu := (d.buyOrder.Price*tax) - d.sellOrder.Price
    out := float64(amt)*ppu

    return out
}

func makeDeal(item *items.Item, bOrder *order.Order, sOrder *order.Order, cDeals chan Deal) bool {
    d := Deal{item, bOrder, sOrder}
    if d.Profit() > 0 {
        deals = append(deals, d)
        cDeals <- d
        return true
    }
    return false
}

func computeBuyOrder(item *items.Item, bOrder *order.Order, sAvl avl.Avl, cDeals chan Deal) {
    iter := sAvl.GetIterator()

    for iter.Next() {
        sOrder := (*iter.Value()).(items.OrderAvlData).Order
        if !makeDeal(item, bOrder, sOrder, cDeals) {
            break
        }
    }
}

func computeSellOrder(item *items.Item, sOrder *order.Order, bAvl avl.Avl, cDeals chan Deal) {
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
