package order

import (
    "time"
    "../utils/avl"
)

type Order struct {
    Duration int `json:"duration"`
    IsBuyOrder bool `json:"is_buy_order"`
    Issued time.Time `json:"issued"`
    LocationID int64 `json:"location_id"`
    MinVolume int `json:"min_volume"`
    OrderID int64 `json:"order_id"`
    Price float64 `json:"price"`
    Range string `json:"range"`
    SystemID int `json:"system_id"`
    ItemID int `json:"type_id"`
    VolumeRemain int `json:"volume_remain"`
    VolumeTotal int `json:"volume_total"`
    Executed int
}

func (a Order) Less (b *avl.Data) bool{
    c := (*b)
    d := c.(Order)
    return a.Price < d.Price
}
