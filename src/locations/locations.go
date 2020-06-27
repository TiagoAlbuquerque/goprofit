package locations

import (
    "fmt"
    "math"
    "io/ioutil"
    "encoding/json"

    "../utils"
)

type location struct{
    Name string `json:"name"`
    Distances map[int64]int `json:"distances"`
}

const fName = "data_locations.eve"
const locationsURL = "https://esi.evetech.net/latest/universe/stations/%d"
const structuresURL = "https://esi.evetech.net/latest/universe/systems/%d"
const systemsURL = "https://esi.evetech.net/latest/universe/systems/%d"

var locations map[int64]location
var saveToFileFlag = false

func getLocationInfo(id int64) location {
    println()
    println("new Location")
    var url string
    var loc location
    println(id)
    if id < 60000000 { 
        url = fmt.Sprintf(systemsURL, id)
    } else if id > 99999999 {
        url = fmt.Sprintf(structuresURL, id)
    } else {
        url = fmt.Sprintf(locationsURL, id)
    }
    println(url)
    utils.JsonFromUrl(url, &loc)
    loc.Distances = map[int64]int{}
    println(loc.Name)
    locations[id] = loc
    saveToFileFlag = true
    return loc
}

func getLocation(locID int64) *location{
    loc, ok := locations[locID]
    if !ok {
        loc = getLocationInfo(locID)
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

func GetDistance(id1, id2 int64) int {
    a := int64(math.Min(float64(id1), float64(id2)))
    b := int64(math.Max(float64(id1), float64(id2)))
    loc := getLocation(a)
    dist, ok := loc.Distances[b]
    if !ok {
        var route []int
        url := fmt.Sprintf("https://esi.evetech.net/latest/route/%d/%d/", a, b)
        //println(url)
        utils.JsonFromUrl(url, &route)        
        loc.Distances[b] = len(route)
        saveToFileFlag = true
    }
    dist, ok = loc.Distances[b]
    return dist
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
