package locations

import (
    "fmt"
    "io/ioutil"
    "encoding/json"

    "../utils"
)

type location struct{
    Name string `json:"name"`
}

const fName = "data_locations.eve"
const locationsURL = "https://esi.evetech.net/latest/universe/stations/%d"
const structuresURL = "https://stop.hammerti.me.uk/api/citadel/%d"

var locations map[int64]location
var saveToFileFlag = false

func getLocationInfo(id int64) location {
    println("new Location")
    var url string
    var loc location
    if id > 99999999 {
        url = fmt.Sprintf(structuresURL, id)
        mloc := map[int64]location{}
        println(url)
        utils.JsonFromUrl(url, &mloc)
        loc = mloc[id]
    } else {
        url = fmt.Sprintf(locationsURL, id)
        println(url)
        utils.JsonFromUrl(url, &loc)
    }
    fmt.Println(loc.Name)
    locations[id] = loc
    saveToFileFlag = true
    return loc
}

func getLocation(locId int64) *location{
    loc, ok := locations[locId]
    if !ok {
        loc = getLocationInfo(locId)
    }
    return &loc
}

func init() {
    raw, err := ioutil.ReadFile(fName)
    if err == nil {
        json.Unmarshal(raw, &locations)
    } else {
        fmt.Printf("Failed to open %s\n", fName)
        locations = make(map[int64]location)
    }
}

func Name(id int64) string{
    loc := getLocation(id)
    return loc.Name
}

// Terminate locations
func Terminate() {
    if !saveToFileFlag { return }
    utils.Save(fName, locations)
    saveToFileFlag = false
}
