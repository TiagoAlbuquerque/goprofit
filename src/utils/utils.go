package utils

import (
    "os"
    "fmt"
    "strings"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "gopkg.in/cheggaaa/pb.v1"
 //   "sort"
)

func GetUrl(url string) *http.Response {
    var res *http.Response
    var err error
    for ok := false; !ok; {
        res, err = http.Get(url)
        ok = (err == nil)
        err = nil
    }
    return res
}

func JsonFromUrl(url string, out interface{}){
    var body []byte
    var err error
    for ok := false; !ok; {
        res := GetUrl(url)
        defer res.Body.Close()
        body, err = ioutil.ReadAll(res.Body)
        ok = (err == nil)
        if !ok { fmt.Print("!OK")}
        err = nil
    }
    json.Unmarshal(body, out)
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
    fmt.Println(f_name, " saved")
}

func StatusIndicator(text string) {
    fmt.Printf("\033[A\r\033[K%s\n", text)
}

func ProgressBar(total int, c chan bool){
    bar := pb.StartNew(total)
    for i:=0; i < total; i++ {
        _ = <-c
        bar.Increment()
    }
    bar.Finish()
}

func commas(s string) string {
    if len(s) <= 3 {
        return s
    } else {
        return commas(s[0:len(s)-3]) + " " + s[len(s)-3:]
    }
}

func FormatCommas(num float64) string {
    parts := strings.Split(fmt.Sprintf("%.2f", num), ".")
    if parts[0][0] == '-' {
        return "-" + commas(parts[0][1:]) + "." + parts[1]
    }
    return commas(parts[0]) + "." + parts[1]
}
