package conf

import(
	"fmt"
	"sync"
	"../utils"
	"io/ioutil"
	"encoding/json"
)

type gpconf struct {
	Cargo float64 `json:"cargo"`
	FItems string `json:"f_items"`
	FLocations string `json:"f_locations"`
	FRegions string `json:"f_regions"`
	MaxInvest float64 `json:"max_invest"`
	MaxTrades int `json:"max_trades"`
	Minpm3 float64 `json:"min_pm3"`
	Tax float64 `json:"tax"`
}

const fname = "data_conf.eve"

var conf gpconf
var saveToFileFlag bool = false
var mutex sync.Mutex

//Cargo availabe in ship
func Cargo() float64{
	return conf.Cargo
}

//Tax value charged in stations
func Tax() float64{
	return conf.Tax
}

//MaxInvest is the maximum amount available to invest in a shopping list
func MaxInvest() float64{
	return conf.MaxInvest
}

//Minpm3 Minimal expected profit amount pem cubic meter of cargo
func Minpm3() float64 {
	return conf.Minpm3
}
 
func init(){
    mutex = sync.Mutex{}
    raw, err := ioutil.ReadFile(fname)
    if err == nil {
        json.Unmarshal(raw, &conf)
    } else {
        fmt.Printf("Failed to open %s\n", fname)
        conf = gpconf{ 100.0, "data_items.eve", "data_locations.eve", "data_regions.eve", 100000000.0, 20, 100000, 0.05}
    }
}

//Terminate method will save possible changes to the configuration file
func Terminate() {
    if !saveToFileFlag { return }

    utils.Save(fname, conf)
    saveToFileFlag = false
}