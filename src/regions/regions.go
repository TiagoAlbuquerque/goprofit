package regions

import (
    "fmt"
    "net/http"
    "encoding/json"
    "io/ioutil"
    "strconv"
//    "os"
    "../utils"
//    "reflect"
//    "math"
)

const regions_url = "https://esi.evetech.net/latest/universe/regions/%s"
const markets_url = "https://esi.evetech.net/latest/markets/%s/orders/?order_type=all&page=%d"
const f_name = "data_regions.eve"

var regions map[string]interface{}

func get_url(url string) *http.Response {

    res, err := http.Get(url)
    for err != nil {
        fmt.Println("ERRO")
        fmt.Println(url)
        fmt.Println(err)
        fmt.Println("tentando novamente")
        err = nil
        res, err = http.Get(url)
    }
    return res
}
func json_from_url(url string) interface{}{
    res := get_url(url)
    defer res.Body.Close()
    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        fmt.Println("couldnt read response body")
        panic(err)
    }
    var out interface{}
    json.Unmarshal(body, &out)
    return out
}

func get_regions_list() []string {

    url := fmt.Sprintf(regions_url, "")
    list := json_from_url(url).([]interface{})
    var out []string
    for i := 0 ; i < len(list); i++ {
        id := fmt.Sprintf("%7.0f", list[i].(float64))
        out = append(out, id)
    }
   return out
}

func get_region_info(id string, c chan bool) {

    url := fmt.Sprintf(regions_url, id)
    info := json_from_url(url).(map[string]interface{})
    info["marked"] = true
    regions[id] = info
    c <- true
}

func get_regions_info(){
    fmt.Println("Getting regions info")
    list := get_regions_list()

    c := make(chan bool)
    total := len(list)
    for i := 0; i < total; i++ {
        go get_region_info(list[i], c)
    }
    utils.ProgressBar(total, c)
}

func getMarketsPagesList() []string {
    var out []string
    for id, i_info := range regions {
        info := i_info.(map[string]interface{})
        if info["marked"].(bool){
            for i:=1; i < info["pages"].(int)+1; i++ {
                url := fmt.Sprintf(markets_url, id, i)
                out = append(out, url)
            }
        }
    }
    return out
}

func get_market_pages_count(id string, c chan bool){
    reg := regions[id].(map[string]interface{})
    url := fmt.Sprintf(markets_url, id, 1)
    var pages []string
    for ok := false; !ok; {
        res := get_url(url)
        defer res.Body.Close()
        pages, ok = res.Header["X-Pages"]
    }
    reg["pages"], _ = strconv.Atoi(pages[0])
    c <- true
}

func update_markets_pages_count(){
    fmt.Println("Updating markets pages count")
    c := make(chan bool)
    total := 0
    for id, reg := range regions {
        if reg.(map[string]interface{})["marked"].(bool) {
            total++
            go get_market_pages_count(id, c)
        }
    }
    utils.ProgressBar(total, c)
}

func getMarketPages(url string, c chan bool){
    var page []interface{}
    for ok := false; !ok; {
        page, ok = json_from_url(url).([]interface{})
    }
    page = page
    c <- true
}
func GetMarketsPages(){
    fmt.Println("Fetching markets pages")
    l := getMarketsPagesList()
    c := make(chan bool)
    total := len(l)
    for i := 0; i < total; i++ {
        go getMarketPages(l[i], c)
    }
    utils.ProgressBar(total, c)
}

func Start(){

    fmt.Println("")
    GetMarketsPages()
}


func init(){
    i_regions, err := utils.Load(f_name)
    if err == nil {
        regions = i_regions.(map[string]interface{})
    } else {
        fmt.Printf("Failed to open %s\n", f_name)
        regions = make(map[string]interface{})
        get_regions_info()
    }

    update_markets_pages_count()
    utils.Save(f_name, regions)
}
