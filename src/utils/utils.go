package utils
import (
    "os"
    "encoding/json"
    "io/ioutil"
       )

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
