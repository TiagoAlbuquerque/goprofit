package items

import(
    "../orders"
    "../utils"
    "../utils/avl"

    "fmt"
    "encoding/json"
    "io/ioutil"
//    "os"
//    "container/list"
    "sync"
//        "strings"
//        "reflect"
)


type Item struct {
    Capacity float32 `json:"capacity"`
    Description string `json:"description"`
    DogmaAttributes []struct {
        AttributeID int `json:"attribute_id"`
        Value float32 `json:"value"`
    } `json:"dogma_attributes"`
    DogmaEffects []struct {
        EffectID int `json:"effect_id"`
        IsDefault bool `json:"is_default"`
    } `json:"dogma_effects"`
    GraphicID int `json:"graphic_id"`
    GroupID int `json:"group_id"`
    MarketGroupID int `json:"market_group_id"`
    Mass float32 `json:"mass"`
    Name string `json:"name"`
    PackagedVolume float32 `json:"packaged_volume"`
    PortionSize int `json:"portion_size"`
    Published bool `json:"published"`
    Radius float32 `json:"radius"`
    ItemID int `json:"type_id"`
    Volume float64 `json:"volume"`

    //Buy_orders *avl.Avl
    //Sell_orders *avl.Avl

    Buy_orders map[int64]int64
    Sell_orders map[int64]int64
}

type OrderAvlData struct {
    Order *orders.Order
}

func (a OrderAvlData) Less (b *avl.Data) bool{
    c := (*b)
    d := c.(OrderAvlData)
    return a.Order.Price < d.Order.Price
}

const itemUrl = "https://esi.evetech.net/latest/universe/types/%d"
const f_name = "data_items.eve"

var items map[int]Item
var saveToFileFlag bool = false
var mutex sync.Mutex

func getItemInfo(id int) Item {
    println("new Item")
    url := fmt.Sprintf(itemUrl, id)
    println(url)
    var item Item
    utils.JsonFromUrl(url, &item)
    fmt.Println(item.Name)
//    item.Buy_orders = avl.NewAvl(avl.REVERSED)
//    item.Sell_orders = avl.NewAvl(avl.DIRECT)
    item.Buy_orders = make(map[int64]int64)
    item.Sell_orders = make(map[int64]int64)
    items[id] = item
    saveToFileFlag = true
    return item
}

func Get(itemId int) Item{
    mutex.Lock()
    defer mutex.Unlock()
    item, ok := items[itemId]
    if !ok {
        item = getItemInfo(itemId)
    }
    return item
}

func Set(item Item) {
    mutex.Lock()
    defer mutex.Unlock()
    items[item.ItemID] = item
}

func (item *Item) place(o orders.Order) {
    if o.IsBuyOrder {
        item.Buy_orders[o.OrderID] = o.OrderID
    } else {
        item.Sell_orders[o.OrderID] = o.OrderID
    }
}

func (item *Item) isOfficer() bool{
    for _, v := range item.DogmaAttributes {
        if v.AttributeID == 1692 && v.Value == 5.0 { return true }
    }
    return false
}

func PlaceOrder(o orders.Order) {
    item := Get(o.ItemID)
    if item.isOfficer() { return }
    item.place(o)
    Set(item)
}

func Cleanup(){
    for _, item := range items {
        //item.Buy_orders = avl.NewAvl(avl.REVERSED)
        //item.Sell_orders = avl.NewAvl(avl.DIRECT)
        item.Buy_orders = make(map[int64]int64)
        item.Sell_orders = make(map[int64]int64)
        Set(item)
    }
}

func init(){
    mutex = sync.Mutex{}
    raw, err := ioutil.ReadFile(f_name)
    if err == nil {
        json.Unmarshal(raw, &items)
    } else {
        fmt.Printf("Failed to open %s\n", f_name)
        items = make(map[int]Item)
    }

    Cleanup()
}

func Terminate() {
    if !saveToFileFlag { return }

    utils.Save(f_name, items)
    saveToFileFlag = false
}
