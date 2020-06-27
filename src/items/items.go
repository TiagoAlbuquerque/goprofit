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

    //BuyOrders *avl.Avl
    //SellOrders *avl.Avl

    BuyOrders map[int64]int64
    SellOrders map[int64]int64
}

type OrderAvlData struct {
    Order *orders.Order
}

func (a OrderAvlData) Less (b *avl.Data) bool{
    c := (*b)
    d := c.(OrderAvlData)
    return a.Order.Price < d.Order.Price
}

const itemURL = "https://esi.evetech.net/latest/universe/types/%d"
const fileName = "data_items.eve"

var items map[int]Item
var saveToFileFlag bool = false
var mutex sync.Mutex

func getItemInfo(id int) Item {
    println()
    println("new Item")
    url := fmt.Sprintf(itemURL, id)
    println(url)
    var item Item
    utils.JsonFromUrl(url, &item)
    println(item.Name)
//    item.BuyOrder = avl.NewAvl(avl.REVERSED)
//    item.SellOrders = avl.NewAvl(avl.DIRECT)
    item.BuyOrders = make(map[int64]int64)
    item.SellOrders = make(map[int64]int64)
    items[id] = item
    saveToFileFlag = true
    return item
}

func Get(itemID int) Item{
    mutex.Lock()
    defer mutex.Unlock()
    item, ok := items[itemID]
    if !ok {
        item = getItemInfo(itemID)
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
        item.BuyOrders[o.OrderID] = o.OrderID
    } else {
        item.SellOrders[o.OrderID] = o.OrderID
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
        //item.BuyOrders = avl.NewAvl(avl.REVERSED)
        //item.SellOrders = avl.NewAvl(avl.DIRECT)
        item.BuyOrders = make(map[int64]int64)
        item.SellOrders = make(map[int64]int64)
        Set(item)
    }
}

func init(){
    mutex = sync.Mutex{}
    raw, err := ioutil.ReadFile(fileName)
    if err == nil {
        json.Unmarshal(raw, &items)
    } else {
        fmt.Printf("Failed to open %s\n", fileName)
        items = make(map[int]Item)
    }

    Cleanup()
}

func Terminate() {
    if !saveToFileFlag { return }

    utils.Save(fileName, items)
    saveToFileFlag = false
}
