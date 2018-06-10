package locations

import (
	"fmt"

	"../utils"
)

const fName = "data_locations.eve"
const locationsURL = "https://esi.evetech.net/latest/universe/stations/%s"
const structuresURL = "https://stop.hammerti.me.uk/api/citadel/%s"

var locations map[string]interface{}
var saveToFileFlag = false

func getLocationInfo(id string) map[string]interface{} {
	var url string
	if len(id) == 8 {
		url = fmt.Sprintf(locationsURL, id)
	} else {
		url = fmt.Sprintf(structuresURL, id)
	}
	fmt.Println(url)
	var location map[string]interface{}
	utils.JsonFromUrl(url, &location)
	locations[id] = location
	saveToFileFlag = true
	return location
}

func init() {
	iLocations, err := utils.Load(fName)
	if err == nil {
		locations = iLocations.(map[string]interface{})
	} else {
		fmt.Printf("Failed to open %s\n", fName)
		locations = make(map[string]interface{})
	}
}

// Terminate locations
func Terminate() {
	if !saveToFileFlag {
		return
	}
	utils.Save(fName, locations)
}
