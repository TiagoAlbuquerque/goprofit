package regions

import (
    "fmt"
    "net/http"
    "encoding/json"
    "io/ioutil"
//    "math"
)

const regions_url = "https://esi.evetech.net/latest/universe/regions/%s"
const markets_url = "https://esi.evetech.net/latest/markets/%s/orders/"

var regions map[string]interface{}

func get_region_info(url string, c chan int){
    fmt.Println(url)
    c <- 1
}

func Start(){
    if regions == nil {
        fmt.Println("no regions")
        res, err := http.Get(fmt.Sprintf(regions_url,""))
        if err != nil {
            fmt.Println("ERRO")
            fmt.Println(err)
        }
        defer res.Body.Close()
        body, err := ioutil.ReadAll(res.Body)

        var list []interface{}
        json.Unmarshal(body, &list)
        c := make(chan int)
        for i := 0 ; i < len(list); i++ {
            id := fmt.Sprintf("%7.0f", list[i].(float64))
//            fmt.Println(fmt.Sprintf(regions_url, id))
            go get_region_info(fmt.Sprintf(regions_url, id), c)
        }

        for count := 0; count < len(list); count++ {
            _ = <-c
            fmt.Println(count)
        }
    }

}

func init(){
    //load from file
//    if err != nil {


//    }
}
