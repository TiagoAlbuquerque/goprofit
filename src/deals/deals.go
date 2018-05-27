package deals

import (
//    "fmt"
    "../items"
    "../order"
//    "../utils/avl"
)

type Deal struct{
    item *items.Item
    buyOrder, sellOrder *order.Order
}


var deals []Deal

func makeDeal(item *items.Item, bOrder *order.Order, sOrder *order.Order) {
    deals = append(deals, Deal{item, bOrder, sOrder})
}

func ComputeDeals(o *order.Order) {
    itID := o.ItemID
    item := items.GetItem(itID)
    item = item

}
