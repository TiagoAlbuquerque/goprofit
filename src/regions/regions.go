package regions

import (
	"../utils"
	//    "../items"
	//    "../order"
	//    "../locations"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/ti/nasync"
	//    "os"
	//    "reflect"
	//    "math"
)

const regionsURL = "https://esi.evetech.net/latest/universe/regions/"
const marketsURL = "https://esi.evetech.net/latest/markets/%d/orders/?order_type=all&page=%d"
const fName = "data_regions.json"

type region struct {
	Constellations []int  `json:"constellations"`
	Description    string `json:"description"`
	Name           string `json:"name"`
	RegionID       int    `json:"region_id"`
	Marked         bool   `json:"marked"`
	pages          int
}

//type mOrder order.Order

var regions map[int]region

var saveToFileFlag bool = false

func getRegionsList() []int {
	var out []int
	utils.JSONFromURL(regionsURL, &out)
	return out
}

func getRegionInfo(id int, c chan bool) {
	url := fmt.Sprint(regionsURL, id)
	var reg region
	utils.JSONFromURL(url, &reg)
	reg.Marked = true
	regions[id] = reg
	saveToFileFlag = true
	c <- true
}

func getRegionsInfo() {
	fmt.Println("Getting regions info")
	list := getRegionsList()
	c := make(chan bool)
	total := len(list)
	async := nasync.New(100, 100)
	defer async.Close()
	for i := 0; i < total; i++ {
		async.Do(getRegionInfo, list[i], c)
	}
	utils.ProgressBar(total, c)
}

//GetMarketsPagesList will produce a list of URLs for market pages
func GetMarketsPagesList() []string {
	var out []string
	for id, reg := range regions {
		if reg.Marked {
			for i := 1; i < reg.pages+1; i++ {
				url := fmt.Sprintf(marketsURL, id, i)
				out = append(out, url)
			}
		}
	}
	return out
}

func getMarketPagesCount(id int, c chan bool) {
	url := fmt.Sprintf(marketsURL, id, 1)
	var pages []string
	for ok := false; !ok; {
		res := utils.GetURL(url)
		defer res.Body.Close()
		pages, ok = res.Header["X-Pages"]
	}
	reg := regions[id]
	reg.pages, _ = strconv.Atoi(pages[0])
	regions[id] = reg
	c <- true
}

func updateMarketsPagesCount() {
	fmt.Println("Updating markets pages count")
	c := make(chan bool)
	total := 0
	async := nasync.New(100, 100)
	defer async.Close()
	for id, reg := range regions {
		if reg.Marked {
			total++
			async.Do(getMarketPagesCount, id, c)
		}
	}
	utils.ProgressBar(total, c)
}

func backup() bool {
	fmt.Printf("Failed to open %s\n", fName)
	regions = make(map[int]region)
	getRegionsInfo()
	return true
}

func init() {
	raw, err := ioutil.ReadFile(fName)
	_ = (err == nil && json.Unmarshal(raw, &regions) == nil) || backup()
	updateMarketsPagesCount()
}

//Terminate will save modifications to the regions to its configuration file
func Terminate() {
	saveToFileFlag = saveToFileFlag && utils.Save(fName, regions) && !saveToFileFlag
}
