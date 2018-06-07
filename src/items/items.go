package items

import(
    "../order"
    "../utils"
    "../utils/avl"

    "fmt"
    "encoding/json"
    "io/ioutil"
//    "os"
//    "container/list"
//        "sync"
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
    Volume float32 `json:"volume"`

    Buy_orders avl.Avl
    Sell_orders avl.Avl
}


const itemUrl = "https://esi.evetech.net/latest/universe/types/%d"
const f_name = "data_items.eve"

var items map[int]Item
var saveToFileFlag bool = false

func getItemInfo(id int) Item {
    println("new Item")
    url := fmt.Sprintf(itemUrl, id)
    println(url)
    var item Item
    utils.JsonFromUrl(url, &item)
    fmt.Println(item.Name)
    item.Buy_orders = avl.Avl{}
    item.Sell_orders = avl.Avl{}
    items[id] = item
    saveToFileFlag = true
    return item
}


func GetItem(itemId int) *Item{
    item, ok := items[itemId]
    if !ok {
        item = getItemInfo(itemId)
    }
    return &item
}

func (item *Item) place(o *order.Order) {
    a := (*o)
    b := avl.Data(a)
    if o.IsBuyOrder {
        item.Buy_orders.Put(&b)
    } else {
        item.Sell_orders.Put(&b)
    }
}

func PlaceOrder(o *order.Order) {
    item := GetItem(o.ItemID)
    item.place(o)
    items[(*item).ItemID] = (*item)
}

func Cleanup(){
    for _, item := range items {
        item.Buy_orders = avl.Avl{}
        item.Sell_orders = avl.Avl{}
    }
}

func init(){
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
    saToFileFlag = false
}
