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
    if !saveToFileFlag { return }
    utils.Save(f_name, locations)
}
