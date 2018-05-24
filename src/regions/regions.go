package regions

import (
    "../utils"
//    "../items"
//    "../order"
//    "../locations"

    "fmt"
    "strconv"
    "encoding/json"
    "io/ioutil"


    "github.com/ti/nasync"
//    "os"
//    "reflect"
//    "math"
        )

const regions_url = "https://esi.evetech.net/latest/universe/regions/"
const markets_url = "https://esi.evetech.net/latest/markets/%d/orders/?order_type=all&page=%d"
const f_name = "data_regions.eve"

type Region struct {
    Constellations []int `json:"constellations"`
    Description string `json:"description"`
    Name string `json:"name"`
    RegionID int `json:"region_id"`
    Marked bool `json:"marked"`
    pages int
}
//type mOrder order.Order

var regions map[int]Region

var saveToFileFlag bool = false

func get_regions_list() []int {
    var out []int
    utils.JsonFromUrl(regions_url, &out)
    return out
}

func get_region_info(id int, c chan bool) {
    url := fmt.Sprint(regions_url, id)
    var region Region
    utils.JsonFromUrl(url, &region)
    region.Marked = true
    regions[id] = region
    saveToFileFlag = true
    c <- true
}

func get_regions_info(){
    fmt.Println("Getting regions info")
    list := get_regions_list()
    c := make(chan bool)
    total := len(list)
    async := nasync.New(100, 100)
    defer async.Close()
    for i := 0; i < total; i++ {
        async.Do(get_region_info, list[i], c)
    }
    utils.ProgressBar(total, c)
}

func GetMarketsPagesList() []string {
    var out []string
    for id, reg := range regions {
        if reg.Marked {
            for i:=1; i < reg.pages+1; i++ {
                url := fmt.Sprintf(markets_url, id, i)
                out = append(out, url)
            }
        }
    }
    return out
}

func get_market_pages_count(id int, c chan bool){
    url := fmt.Sprintf(markets_url, id, 1)
    var pages []string
    for ok := false; !ok; {
        res := utils.GetUrl(url)
        defer res.Body.Close()
        pages, ok = res.Header["X-Pages"]
    }
    reg := regions[id]
    reg.pages, _ = strconv.Atoi(pages[0])
    regions[id] = reg
    c <- true
}

func update_markets_pages_count(){
    fmt.Println("Updating markets pages count")
    c := make(chan bool)
    total := 0
    async := nasync.New(100, 100)
    defer async.Close()
    for id, reg := range regions {
        if reg.Marked {
            total++
            async.Do(get_market_pages_count, id, c)
        }
    }
    utils.ProgressBar(total, c)
}

/*
func getMarketPages(url string, cPages chan []order.Order){
    var mPage []order.Order
    for ok := false; !ok; {
        res := utils.GetUrl(url)
        defer res.Body.Close()
        body, _ := ioutil.ReadAll(res.Body)
        json.Unmarshal(body, &mPage)
        ok = (mPage != nil)
    }
    cPages <- mPage
}

func GetMarketsPages(){
    fmt.Println("Fetching markets pages")
    lURL := GetMarketsPagesList()
    cOK := make(chan bool)
    cPages := make(chan []order.Order)
    total := len(lURL)
    go items.ConsumePages(cPages, cOK, total)
    async := nasync.New(1000, 1000)
    defer async.Close()
    for i := 0; i < total; i++ {
        async.Do(getMarketPages, lURL[i], cPages)
    }
    utils.ProgressBar(total, cOK)
}
*/

func init(){
    raw, err := ioutil.ReadFile(f_name)
    if err == nil {
        json.Unmarshal(raw, &regions)
    } else {
        fmt.Printf("Failed to open %s\n", f_name)
        regions = make(map[int]Region)
        get_regions_info()
    }

    update_markets_pages_count()
}

func Terminate(){
    if !saveToFileFlag {
        return
    }
    utils.Save(f_name, regions)
}
