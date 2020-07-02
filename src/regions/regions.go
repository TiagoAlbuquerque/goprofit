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
const fName = "data_regions.eve"

type Region struct {
	Constellations []int  `json:"constellations"`
	Description    string `json:"description"`
	Name           string `json:"name"`
	RegionID       int    `json:"region_id"`
	Marked         bool   `json:"marked"`
	pages          int
}

//type mOrder order.Order

var regions map[int]Region

var saveToFileFlag bool = false

func getRegionsList() []int {
	var out []int
	utils.JsonFromUrl(regionsURL, &out)
	return out
}

func getRegionInfo(id int, c chan bool) {
	url := fmt.Sprint(regionsURL, id)
	var region Region
	utils.JsonFromUrl(url, &region)
	region.Marked = true
	regions[id] = region
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

func init() {
	raw, err := ioutil.ReadFile(fName)
	if err == nil {
		json.Unmarshal(raw, &regions)
	} else {
		fmt.Printf("Failed to open %s\n", fName)
		regions = make(map[int]Region)
		getRegionsInfo()
	}

	updateMarketsPagesCount()
}

func Terminate() {
	if !saveToFileFlag {
		return
	}
	utils.Save(fName, regions)
	saveToFileFlag = false
}
