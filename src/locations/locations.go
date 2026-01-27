package locations

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"sync"

	"goprofit/utils"
)

type location struct {
	Name      string        `json:"name"`
	Distances map[int64]int `json:"distances"`
}

const fName = "data_locations.json"
const locationsURL = "https://esi.evetech.net/latest/universe/stations/%d"
const structuresURL = "https://esi.evetech.net/latest/universe/systems/%d"
const systemsURL = "https://esi.evetech.net/latest/universe/systems/%d"

var locations map[int64]location
var saveToFileFlag = false
var mutex sync.Mutex // Only protects write operations

func getLocationInfo(id int64) location {
	var url string
	var loc location
	if id < 60000000 {
		url = fmt.Sprintf(systemsURL, id)
	} else if id > 99999999 {
		url = fmt.Sprintf(structuresURL, id)
	} else {
		url = fmt.Sprintf(locationsURL, id)
	}
	utils.StatusLine(15, "new location: "+url)
	utils.JSONFromURL(url, &loc)
	loc.Distances = map[int64]int{}
	utils.StatusLine(15, loc.Name)
	// Note: caller must hold write lock
	locations[id] = loc
	saveToFileFlag = true
	return loc
}

func getLocation(locID int64) *location {
	// Lock-free read - locations data is stable once written
	loc, ok := locations[locID]
	if !ok {
		// Fetch without lock to avoid blocking other goroutines
		var newLoc location
		var url string
		if locID < 60000000 {
			url = fmt.Sprintf(systemsURL, locID)
		} else if locID > 99999999 {
			url = fmt.Sprintf(structuresURL, locID)
		} else {
			url = fmt.Sprintf(locationsURL, locID)
		}
		utils.StatusLine(15, "new location: "+url)
		utils.JSONFromURL(url, &newLoc)
		newLoc.Distances = map[int64]int{}
		utils.StatusLine(15, newLoc.Name)

		// Only lock for write
		mutex.Lock()
		locations[locID] = newLoc
		saveToFileFlag = true
		mutex.Unlock()
		return &newLoc
	}
	return &loc
}

// GetDistance will return the number of jumps on a route from id1 to id2
func GetDistance(id1, id2 int64) int {
	a := int64(math.Min(float64(id1), float64(id2)))
	b := int64(math.Max(float64(id1), float64(id2)))
	loc := getLocation(a)

	// Lock-free read of existing distance
	dist := loc.Distances[b]
	if dist != 0 {
		return dist
	}

	// Need to fetch route
	var route []int
	url := fmt.Sprintf("https://esi.evetech.net/latest/route/%d/%d/", a, b)
	utils.StatusLine(15, url)
	utils.JSONFromURL(url, &route)

	newDist := len(route)
	if newDist == 0 {
		newDist = 1
	}

	// Only lock for write
	mutex.Lock()
	loc.Distances[b] = newDist
	saveToFileFlag = true
	mutex.Unlock()

	return newDist
}

// GetName will return the name of the system for a defined id
func GetName(id int64) string {
	loc := getLocation(id)
	return loc.Name
}

func backup() bool {
	fmt.Printf("Failed to open %s\n", fName)
	locations = make(map[int64]location)
	return true
}

func init() {
	raw, err := ioutil.ReadFile(fName)
	_ = (err == nil && json.Unmarshal(raw, &locations) == nil) || backup()
}

// Terminate locations
func Terminate() {
	mutex.Lock()
	defer mutex.Unlock()
	if saveToFileFlag {
		err := utils.SaveToJSONFile(fName, locations)
		if err != nil {
			fmt.Println("Failed to save locations to file:", err)
		}
		saveToFileFlag = false
	}
}
