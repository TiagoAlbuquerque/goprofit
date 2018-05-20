package utils
import (
    "os"
    "encoding/json"
    "io/ioutil"
    "gopkg.in/cheggaaa/pb.v1"
    "net/http"
    "fmt"
    "sort"
)

func GetUrl(url string) *http.Response {
    var res *http.Response
    var err error
    for ok := false; !ok; {
        res, err = http.Get(url)
        ok = (err == nil)
    }
    return res
}

func InsertSorted(l []interface{}, o map[string]interface{}, reversed bool ) []interface{}{
    mult := 1.0;
    if reversed { mult = -1.0 }
    comp := func (i int) bool {
        a := mult*((l[i].(map[string]interface{}))["price"].(float64))
        b := mult*(o["price"].(float64))
        return (a > b)
    }
    i := sort.Search(len(l), comp)
    l = append(l, make(map[string]interface{}))
    copy(l[i+1:], l[i:])
    l[i] = o
    return l
}

func JsonFromUrl(url string) interface{}{
    var out interface{}
    var body []byte
    var err error
    for ok := false; !ok; {
        res := GetUrl(url)
        defer res.Body.Close()
        body, err = ioutil.ReadAll(res.Body)
        ok = (err == nil)
        if !ok { fmt.Print("!OK")}
    }
    json.Unmarshal(body, &out)
    return out
}

func Load(f_name string) (interface{}, error) {
    raw, err := ioutil.ReadFile(f_name)
    if err != nil {
        return nil, err
    }

    var c interface{}
    json.Unmarshal(raw, &c)
    return c, nil
}

func Save(f_name string, data interface{}){
    out, _ := json.MarshalIndent(data, "", "  ")
    f, err := os.Create(f_name)
    if err != nil {
        panic(err)
    }
    defer f.Close()
    f.Write(out)
}

func ProgressBar(total int, c chan bool){
    bar := pb.StartNew(total)
    for i:=0; i < total; i++ {
        _ = <-c
        bar.Increment()
    }
    bar.Finish()
}

