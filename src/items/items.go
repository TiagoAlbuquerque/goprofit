package items

import(
    "../utils"
    "../deals"
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

func getItemInfo(id string) map[string]interface{} {
    println("new Item")
    url := fmt.Sprintf(itemUrl, id)
    println(url)
    item := utils.JsonFromUrl(url).(map[string]interface{})
    fmt.Println(item["name"])
    item["buy_orders"] = make([]interface{}, 0)
    item["sell_orders"] = make([]interface{}, 0)
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

func place(order map[string]interface{}, item map[string]interface{}) {
    if order["is_buy_order"].(bool){
        item["buy_orders"] = utils.InsertSorted(item["buy_orders"].([]interface{}), order, true)
        deals.ComputeDeals(item, []interface{}{order}, item["sell_orders"].([]interface{}))
    } else {
        item["sell_orders"] = utils.InsertSorted(item["sell_orders"].([]interface{}), order, false)
        deals.ComputeDeals(item, item["buy_orders"].([]interface{}), []interface{}{order})
    }
}
func place1order(order map[string]interface{}) {
    itemId := fmt.Sprint(order["type_id"])
    item := GetItem(itemId)
    place(order, item)

}

func placeOrders(orders []interface{}) {
    for _, i_order := range orders {
        i_order = i_order
        place1order(i_order.(map[string]interface{}))
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
        item["buy_orders"] = make([]interface{}, 0)
        item["sell_orders"] = make([]interface{}, 0)
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
