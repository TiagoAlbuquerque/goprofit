package main

import (
        "bytes"
        "encoding/gob"
        "encoding/json"
        "fmt"
        "io/ioutil"
        "net/http"
        //"os"
       )

const(
        f_restrictions = "restrictions.eve";
     )
var (
        restrictions    map[string]interface{};
        items           map[string]interface{};
        regions         map[string]interface{};
        stations        map[string]interface{};
    )
func read_crest(url string) map[string]interface{} {
    js := json_url(url)
    if js["next"] {
        print("more")
    }
    return js
}
func json_url(url string) map[string]interface{} {
    resp, err := http.Get(url)
    if err != nil {
        fmt.Printf("erro: %v\n", err)
        return nil
    }
    defer resp.Body.Close();
    body, _ := ioutil.ReadAll(resp.Body)
    var res map[string]interface{};
    json.Unmarshal(body, &res);
    return res;
}
func load_file(f_name string, f func()map[string]interface{} ) map[string]interface{} {
    file, e := ioutil.ReadFile(f_name);

    if e != nil {
        fmt.Printf("file error: %v\n", e);
        return f();
    }
    var res map[string]interface{};
    json.Unmarshal(file, &res);
    return res;
}

func load_restrictions(){
    restrictions = load_file(f_restrictions, get_restrictions);
}

func get_bytes(key interface{}) []byte {
    var buf bytes.Buffer;
    enc := gob.NewEncoder(&buf);
    err := enc.Encode(key);
    if err != nil {
        return nil
    }
    return buf.Bytes()
}

func get_restrictions() map[string]interface{}{
    res := map[string]interface{}{"top":10, "tax":0.05}
    fmt.Printf("created basic restrictions: %v", res);
    return res
}

func main(){
    load_restrictions();
    read_crest("https://crest-tq.eveonline.com")
}

