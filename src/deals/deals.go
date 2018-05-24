package deals

import (
//    "fmt"
    "../order"
    "../utils/avl"
)

type Deal struct{
    item map[string]interface{}
    buyOrder, sellOrder order.Order
}


var deals []Deal

func makeDeal(item map[string]interface{}, bOrder order.Order, sOrder order.Order) {
    deals = append(deals, Deal{item, bOrder, sOrder})
}

func ComputeDealsA(item map[string]interface{}, o order.Order, a avl.Avl) {
    if o.IsBuyOrder {
    }
}


func ComputeDeals(item map[string]interface{}, bList []interface{}, sList []interface{}) {
    /*for _, i_buy := range bList{
        buy := i_buy.(map[string]interface{})
        bPrice := buy["price"].(float64)
        for _, i_sell := range sList {
            sell := i_sell.(map[string]interface{})
            sPrice := sell["price"].(float64)
            if sPrice < (bPrice*1.01){ break }
            makeDeal(item, buy, sell)
        }
    }*/
}
