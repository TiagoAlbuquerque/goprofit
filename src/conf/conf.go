package conf

import(
	"fmt"
	"sync"
	"../utils"
	"io/ioutil"
	"encoding/json"
)

type Conf struct {
	Cargo float64 `json:"cargo"`
	F_Items string `json:"f_items"`
	F_Locations string `json:"f_locations"`
	F_Regions string `json:"f_regions"`
	Min_pm3 int `json:"min_pm3"`
	Tax float32 `json:"tax"`
}

const f_name = "data_conf.eve"

var conf Conf
var saveToFileFlag bool = false
var mutex sync.Mutex

func Cargo() float64{
	return conf.Cargo
}

func init(){
    mutex = sync.Mutex{}
    raw, err := ioutil.ReadFile(f_name)
    if err == nil {
        json.Unmarshal(raw, &conf)
    } else {
        fmt.Printf("Failed to open %s\n", f_name)
        conf = Conf{ 0.0, "", "", "", 0, 0.05}
    }

    //Cleanup()
}

func Terminate() {
    if !saveToFileFlag { return }

    utils.Save(f_name, conf)
    saveToFileFlag = false
}