package regions

import (
    "fmt"
    "net/http"
    "encoding/json"
    "io/ioutil"
    "strconv"
    "os"
//    "reflect"
    "gopkg.in/cheggaaa/pb.v1"
//    "math"
)

const regions_url = "https://esi.evetech.net/latest/universe/regions/%s"
const markets_url = "https://esi.evetech.net/latest/markets/%s/orders/"
const f_name = "data_regions.eve"

var regions map[string]interface{}

func get_url(url string) *http.Response {

    res, err := http.Get(url)
    for err != nil {
        fmt.Println("ERRO")
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

func get_region_info(id string, c chan map[string]interface{}) {

    url := fmt.Sprintf(regions_url, id)
    info := json_from_url(url).(map[string]interface{})
    info["marked"] = true
    c <- info
}

func get_regions_info(){
    fmt.Println("Getting regions info")
    list := get_regions_list()

    c := make(chan map[string]interface{})
    total := len(list)
    for i := 0; i < total; i++ {
        go get_region_info(list[i], c)
    }

    bar := pb.StartNew(total)
    bar.ShowElapsedTime = true
    for i := 0; i < total; i++ {
        info := <-c
        var id = fmt.Sprintf("%7.0f", info["region_id"].(float64))
        regions[id] = info
        bar.Increment()
    }
    bar.Finish()
}

func get_market_pages_count(id string, c chan map[string]interface{}){
    url := fmt.Sprintf(markets_url, id)
    out := make(map[string]interface{})
    out["id"] = id
    res := get_url(url)
    defer res.Body.Close()
    pages, _ := strconv.Atoi(res.Header["X-Pages"][0])
    out["pages"] = pages
    c <- out
}


func update_markets_pages_count(){
    fmt.Println("Updating markets pages count")
    c := make(chan map[string]interface{})
    total := 0
    for id, reg := range regions {
        if reg.(map[string]interface{})["marked"].(bool) {
            total++
            go get_market_pages_count(id, c)
        }
    }
    bar := pb.StartNew(total)
    bar.ShowElapsedTime = true
    for i:=0; i < total; i++ {
        m := <-c
        id := m["id"].(string)
        pages := m["pages"].(int)
        regions[id].(map[string]interface{})["pages"] = pages
        bar.Increment()
    }
    bar.Finish()

}
func Start(){

    fmt.Println("")
}

func load() interface{}{
    return 1
    return nil
}

func save(){
    out, _ := json.MarshalIndent(regions, "", "  ")
    f, err := os.Create(f_name)
    if err != nil {
        panic(err)
    }
    defer f.Close()
    f.Write(out)
}

func init(){
    err := load()
    if err != nil {
        fmt.Printf("Failed to open %s\n", f_name)
        regions = make(map[string]interface{})
        get_regions_info()
    }
    update_markets_pages_count()
    save()
}
