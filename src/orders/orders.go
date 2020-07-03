package orders

import (
	"sync"
	"time"
)

//Order mimics the structure of EVE Online ESI market OrdersMap
type Order struct {
	Duration     int       `json:"duration"`
	IsBuyOrder   bool      `json:"is_buy_order"`
	Issued       time.Time `json:"issued"`
	LocationID   int64     `json:"location_id"`
	MinVolume    int       `json:"min_volume"`
	OrderID      int64     `json:"order_id"`
	Price        float64   `json:"price"`
	Range        string    `json:"range"`
	SystemID     int64     `json:"system_id"`
	ItemID       int       `json:"type_id"`
	VolumeRemain int       `json:"volume_remain"`
	VolumeTotal  int       `json:"volume_total"`
	Executed     int
}

var orders map[int64]*Order
var mutex sync.Mutex

//OrderRemain will return how much of the order is still available to be executed
func (o *Order) OrderRemain() int {
	return o.VolumeRemain - o.Executed
}

//Execute will fill up the order on qnt amount
func (o *Order) Execute(qnt int) {
	o.Executed += qnt
	//Set(*o)
}

//Reset will reset the order to an unexecuted state
func (o *Order) Reset() {
	o.Executed = 0
	//Set(*o)
}

//Get will return the market order to an specific ID
func Get(oID int64) *Order {
	mutex.Lock()
	defer mutex.Unlock()
	out := orders[oID]
	return out
}

//Set will store the receiver Order
func Set(o Order) {
	mutex.Lock()
	defer mutex.Unlock()
	orders[o.OrderID] = &o
}

//Cleanup will empty down the OrdersMap list
func Cleanup() {
	orders = make(map[int64]*Order)
}

func init() {
	mutex = sync.Mutex{}
	Cleanup()
}
