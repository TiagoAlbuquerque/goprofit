package deals

import (
//    "fmt"
)

type Deal struct{
    item, buyOrder, sellOrder map[string]interface{}
}


var deals []Deal

func makeDeal(item map[string]interface{}, bOrder map[string]interface{}, sOrder map[string]interface{}) {
    deals = append(deals, Deal{item, bOrder, sOrder})
}

func ComputeDeals(item map[string]interface{}, bList []interface{}, sList []interface{}) {
    for _, i_buy := range bList{
        buy := i_buy.(map[string]interface{})
        bPrice := buy["price"].(float64)
        for _, i_sell := range sList {
            sell := i_sell.(map[string]interface{})
            sPrice := sell["price"].(float64)
            if sPrice < (bPrice*1.01){ break }
            makeDeal(item, buy, sell)
        }
    }
}
