package order

import (
    "fmt"
)

type Order struct{
    data map[string]interface{}
}

func New(data map[string]interface{}) Order {
    return Order{data}
}
func (o *Order) Price() float64 {
    return o.data["price"].(float64)
}

func (o *Order) IsBuyOrder() bool {
    return o.data["is_buy_order"].(bool)
}

func (o *Order) ItemId() string {
    return fmt.Sprintf("%.0f", o.data["type_id"].(float64))
}

func (a *Order) Less(b Order) bool {
    return a.Price() < b.Price()
}
