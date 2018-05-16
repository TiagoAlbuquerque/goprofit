package regions

import (
    "fmt"
    "net/http"
    "encoding/json"
    "io/ioutil"
    "strconv"
//    "math"
)

const regions_url = "https://esi.evetech.net/latest/universe/regions/%s"
const markets_url = "https://esi.evetech.net/latest/markets/%s/orders/"

var regions map[string]interface{}

func get_url(url string) *http.Response {

    res, err := http.Get(url)
    for err != nil {
        fmt.Println("ERRO")
        fmt.Println(err)
        fmt.Println("tentando novamente")
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
    c <- info
}

func get_regions_info(){
    list := get_regions_list()

    c := make(chan map[string]interface{})
    for i := 0; i < len(list); i++ {
        go get_region_info(list[i], c)
    }
    for i := 0; i < len(list); i++ {
        info := <-c
        var id = fmt.Sprintf("%7.0f", info["region_id"].(float64))
        regions[id] = info
    }
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
    c := make(chan map[string]interface{})
    for id, _ := range regions {
        //if reg.(map[string]interface{})["marked"] {
            go get_market_pages_count(id, c)
        //}
    }
    for range regions {
        m := <-c
        id := m["id"].(string)
        pages := m["pages"].(int)
        regions[id].(map[string]interface{})["pages"] = pages
    }

}
func Start(){
    if regions == nil {
        regions = make(map[string]interface{})
        fmt.Println("no regions")
        get_regions_info()
        update_markets_pages_count()
    }

}

func init(){
    //load from file
//    if err != nil {


//    }
    update_markets_pages_count()
}