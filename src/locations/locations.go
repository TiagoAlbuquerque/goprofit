package locations

import (
    "fmt"
    "../utils"
)

const f_name = "data_locations.eve"
const locationsUrl = "https://esi.evetech.net/latest/universe/stations/%s"
const structuresUrl = "https://stop.hammerti.me.uk/api/citadel/%s"
var locations map[string]interface{}
var saveToFileFlag bool = false

func getLocationInfo(id string) map[string]interface{} {
    var url string
    if len(id) == 8 { url = fmt.Sprintf(locationsUrl, id)
    } else { url = fmt.Sprintf(structuresUrl, id) }
    fmt.Println(url)
    var location map[string]interface{}
    utils.JsonFromUrl(url, &location)
    cleanLocation(location)
    locations[id] = location
    saveToFileFlag = true
    return location
}
/*
func getLocationForOrder(order map[string]interface{}) map[string]interface{} {
    locId := fmt.Sprintf("%.0f",order["location_id"])
    location, ok := locations[locId]
    if !ok {
        locations[locId] = make(map[string]interface{})
        cleanLocation(locations[locId].(map[string]interface{}))
        location = locations[locId]
    }
    return location.(map[string]interface{})
}

func place(order map[string]interface{}, location map[string]interface{}) {
    locItems := location["items"].(map[string]interface{})
    buy_orders := locItems["buy_orders"].([]interface{})
    sell_orders := locItems["sell_orders"].([]interface{})
    if order["is_buy_order"].(bool){
        locItems["buy_orders"] = utils.InsertSorted(buy_orders, order, true)
    } else {
        locItems["sell_orders"] = utils.InsertSorted(sell_orders, order, false)
    }
}

func place1order(order map[string]interface{}) {
    location := getLocationForOrder(order)
    place(order, location)
}

func placeOrders(orders []interface{}) {
    for _, order := range orders {
        place1order(order.(map[string]interface{}))
    }
}

func ConsumePages(cPages chan []interface{}, cOK chan bool, total int) {
    cleanup()
    for i := 0; i < total; i++ {
        placeOrders(<-cPages)
        cOK <- true
    }
}

*/
func cleanLocation(location map[string]interface{}) {
        locItems := make(map[string]interface{})
        locItems["buy_orders"] = make([]interface{}, 0)
        locItems["sell_orders"] = make([]interface{}, 0)
        location["items"] = locItems
}
func cleanup() {
    for _, i_location := range locations {
        location := i_location.(map[string]interface{})
        cleanLocation(location)
    }
}

func init(){
    i_locations, err := utils.Load(f_name)
    if err == nil {
        locations = i_locations.(map[string]interface{})
    } else {
        fmt.Printf("Failed to open %s\n", f_name)
        locations = make(map[string]interface{})
    }
    cleanup()
}

func Terminate() {
    if !saveToFileFlag {
        return
    }
    for _, i_location := range locations {
        location := i_location.(map[string]interface{})
        delete(location, "items")
        delete(location, "deals")
    }
    utils.Save(f_name, locations)
}
