package deals

import (
//    "fmt"
    "../items"
    "../order"
    "../utils/avl"
)

type Deal struct{
    item *items.Item
    buyOrder, sellOrder *order.Order
}

var deals []Deal

func (*Deal) pm3() float64 {
    out := 0.0

    return out
}

func (*Deal) amount(availCargo float64) int {
    out := 0

    return out
}

func (*Deal) profit() float64 {
    out := 0.0

    return out
}

func makeDeal(item *items.Item, bOrder *order.Order, sOrder *order.Order) bool {
    d := Deal{item, bOrder, sOrder}
    if d.profit() > 0 {
        deals = append(deals, Deal{item, bOrder, sOrder})

        return true
    }
    return false
}

func computeBuyOrder(item *items.Item, bOrder *order.Order, sAvl avl.Avl) {
    iter := sAvl.GetIterator(false)

    for iter.Next() {
        sOrder := (*iter.Value()).(order.Order)
        if !makeDeal(item, bOrder, &sOrder) {
            break
        }
    }
}

func computeSellOrder(item *items.Item, sOrder *order.Order, bAvl avl.Avl) {
    iter := bAvl.GetIterator(true)

    for iter.Next() {
        bOrder := (*iter.Value()).(order.Order)
        if !makeDeal(item, &bOrder, sOrder) {
            break
        }
    }
}

func ComputeDeals(o *order.Order) {
    itID := o.ItemID
    item := items.GetItem(itID)
    if o.IsBuyOrder {
        computeBuyOrder(item, o, item.Sell_orders)
    } else {
        computeSellOrder(item, o, item.Buy_orders)
    }

}
