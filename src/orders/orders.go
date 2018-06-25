package orders

import (
    "time"
    "sync"
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

var orders map[int64]Order
var mutex sync.Mutex

func (o *Order) OrderRemain() int {
    return o.VolumeRemain - o.Executed
}

func (o *Order) Execute(qnt int) {
    o.Executed += qnt
    Set(*o)
}

func (o *Order) Reset() {
    o.Executed = 0
    Set(*o)
}

func Get(oID int64) Order {
    mutex.Lock()
    defer mutex.Unlock()
    out := orders[oID]
    return out
}

func Set(o Order) {
    mutex.Lock()
    defer mutex.Unlock()
    orders[o.OrderID] = o
}

func Cleanup() {
    orders = make(map[int64]Order)
}

func init() {
    mutex = sync.Mutex{}
    Cleanup()
}
