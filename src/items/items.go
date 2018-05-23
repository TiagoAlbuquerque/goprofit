package items

import(
//    "../deals"
    "../order"
    "../utils"
    "../utils/avl"
    "fmt"
//    "os"
//    "container/list"
//        "sync"
//        "strings"
//        "reflect"
)

const itemUrl = "https://esi.evetech.net/latest/universe/types/%s"
const f_name = "data_items.eve"
var items map[string]interface{}
var saveToFileFlag bool = false

type Item struct {
    data map[string]interface{}
}

type mOrder order.Order

func (a mOrder) Less (b avl.Data) bool{
    return a.Less(b.(mOrder))

}

func getItemInfo(id string) map[string]interface{} {
    println("new Item")
    url := fmt.Sprintf(itemUrl, id)
    println(url)
    item := utils.JsonFromUrl(url).(map[string]interface{})
    fmt.Println(item["name"])
//    item["buy_orders"] = []order.Order{}
//    item["sell_orders"] = []order.Order{}
    item["buy_orders"] = avl.Avl{}
    item["sell_orders"] = avl.Avl{}
    items[id] = item
    saveToFileFlag = true
    return item
}


func GetItem(itemId string) map[string]interface{}{
    item, ok := items[itemId]
    if !ok {
        item = getItemInfo(itemId)
    }
    return item.(map[string]interface{})
}

func place(o order.Order, item map[string]interface{}) {
    if o.IsBuyOrder() {
        (item["buy_orders"].(avl.Avl)).Put(mOrder(o))
//        item["buy_orders"] = utils.InsertSorted(item["buy_orders"].([]order.Order), o, true)
        //deals.ComputeDeals(item, []order.Order{o}, item["sell_orders"].([]interface{}))
    } else {
        (item["sell_orders"].(avl.Avl)).Put(mOrder(o))
//        item["sell_orders"] = utils.InsertSorted(item["sell_orders"].([]order.Order), o, false)
        //deals.ComputeDeals(item, item["buy_orders"].([]interface{}), []order.Order{o})
    }
}
func place1order(o order.Order) {
    itemId := o.ItemId()
    item := GetItem(itemId)
    place(o, item)

}

func placeOrders(orders []interface{}) {
    for _, o := range orders {
        place1order(order.New(o.(map[string]interface{})))
    }
}

func ConsumePages(cPages chan []interface{}, cOK chan bool, total int) {
    cleanup()
    for i := 0; i < total; i++ {
        placeOrders(<-cPages)
        cOK <- true
    }
}

func cleanup(){
    for _, i_item := range items {
        item := i_item.(map[string]interface{})
//        item["buy_orders"] = []order.Order{}
//        item["sell_orders"] = []order.Order{}
        item["buy_orders"] = avl.Avl{}
        item["sell_orders"] = avl.Avl{}
    }
}

func init(){
    i_items, err := utils.Load(f_name)
    if err == nil {
        items = i_items.(map[string]interface{})
    } else {
        fmt.Printf("Failed to open %s\n", f_name)
        items = make(map[string]interface{})
    }
    cleanup()
}

func Terminate() {
    if !saveToFileFlag {
        return
    }
    for _, i_item := range items {
        item := i_item.(map[string]interface{})
        delete(item, "buy_orders")
        delete(item, "sell_orders")
    }
    utils.Save(f_name, items)
}
