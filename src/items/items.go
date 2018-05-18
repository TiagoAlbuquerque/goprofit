package items

import(
        "../utils"
//        "os"
        "fmt"
//        "reflect"
)


const f_name = "data_items.eve"
var items map[string]interface{}

func getItemForOrder(order map[string]interface{}) map[string]interface{}{
    itemId := order["type_id"]
    itemId = itemId
    return nil
}
func addSorted(l []interface{}, o map[string]interface{}, reversed bool){
    mult := 1
    if reversed {
        mult = mult*(-1)
    }
}
func place(order map[string]interface{}, item map[string]interface{}) {
    if order["is_buy_order"].(bool){
        addSorted(item["buy_orders"].([]interface{}), order, true)
    } else {
        addSorted(item["buy_orders"].([]interface{}), order, false)
    }
}
func placeOrder(order map[string]interface{}) {
    item := getItemForOrder(order)
    item = item
    //place(order, item)

}

func PlaceOrders(orders []interface{}) {
    for _, i_order := range orders {
//        order, ok := i_order.(map[string]interface{})
        placeOrder(i_order.(map[string]interface{}))
    }
    //os.Exit(1)
}

func init(){
    i_items, err := utils.Load(f_name)
    if err == nil {
        items = i_items.(map[string]interface{})
    } else {
        fmt.Printf("Failed to open %s\n", f_name)
        items = make(map[string]interface{})
        utils.Save(f_name, items)
    }
}
