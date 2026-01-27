package regions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"sync"

	"goprofit/utils"
)

const (
	regionsURL   = "https://esi.evetech.net/latest/universe/regions/"
	marketsURL   = "https://esi.evetech.net/latest/markets/%d/orders/?order_type=all&page=%d"
	fileName     = "data_regions.json"
	requestLimit = 100
)

type region struct {
	Constellations []int  `json:"constellations"`
	Description    string `json:"description"`
	Name           string `json:"name"`
	RegionID       int    `json:"region_id"`
	Marked         bool   `json:"marked"`
	Pages          int
}

var (
	regions        map[int]region
	regionsMutex   sync.Mutex
	saveToFileFlag bool
)

func getRegionsList() []int {
	var out []int
	utils.JSONFromURL(regionsURL, &out)
	return out
}
func GetRegionsList() map[int]region {
	return regions
}

func getRegionInfo(id int, c chan interface{}) {
	url := fmt.Sprint(regionsURL, id)
	var reg region
	utils.JSONFromURL(url, &reg)
	reg.Marked = true
	regionsMutex.Lock()
	regions[id] = reg
	saveToFileFlag = true
	regionsMutex.Unlock()
	c <- true
}

func getRegionsInfo() {
	fmt.Println("Getting regions info")
	list := getRegionsList()
	c := make(chan interface{})
	total := len(list)
	for i := 0; i < total; i++ {
		go getRegionInfo(list[i], c)
	}
	utils.ProgressBar(total, c)
}

// GetMarketsPagesList will produce a list of URLs for market pages
func GetMarketsPagesList() []string {
	// Lock-free read - regions map is stable after init
	var out []string
	for id, reg := range regions {
		if reg.Marked {
			for i := 1; i < reg.Pages+1; i++ {
				url := fmt.Sprintf(marketsURL, id, i)
				out = append(out, url)
			}
		}
	}
	return out
}

func getMarketPagesCount(id int, c chan interface{}) {
	url := fmt.Sprintf(marketsURL, id, 1)
	var pages []string
	for ok := false; !ok; {
		res, err := utils.GetURL(url)
		if err != nil || res == nil {
			continue
		}
		pages, ok = res.Header["X-Pages"]
		res.Body.Close() // Close immediately, not defer (avoid accumulation in loop)
	}
	regionsMutex.Lock()
	reg := regions[id]
	reg.Pages, _ = strconv.Atoi(pages[0])
	regions[id] = reg
	regionsMutex.Unlock()
	c <- true
}

func updateMarketsPagesCount() {
	fmt.Println("Updating markets pages count")
	c := make(chan interface{})
	total := 0
	for id, reg := range regions {
		if reg.Marked {
			total++
			go getMarketPagesCount(id, c)
		}
	}
	utils.ProgressBar(total, c)
}

func backup() bool {
	fmt.Printf("Failed to open %s\n", fileName)
	regions = make(map[int]region)
	getRegionsInfo()
	return true
}

func init() {
	raw, err := ioutil.ReadFile(fileName)
	_ = (err == nil && json.Unmarshal(raw, &regions) == nil) || backup()
	updateMarketsPagesCount()
}

// Terminate will save modifications to the regions to its configuration file
// Terminate saves the regions to the configuration file.
func Terminate() {
	regionsMutex.Lock()
	defer regionsMutex.Unlock()

	if saveToFileFlag {
		err := utils.SaveToJSONFile(fileName, regions)
		if err != nil {
			fmt.Println("Failed to save regions to file:", err)
		}
		saveToFileFlag = false
	}
}
