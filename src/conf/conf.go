package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"

	"../utils"
)

type gpconf struct {
	Cargo      float64 `json:"cargo"`
	FItems     string  `json:"f_items"`
	FLocations string  `json:"f_locations"`
	FRegions   string  `json:"f_regions"`
	MaxInvest  float64 `json:"max_invest"`
	MThreshold float64 `json:"message_threshold"`
	MPhoneNum  string  `json:"wapp_phone"`
	MaxTrades  int     `json:"max_trades"`
	Minpm3     float64 `json:"min_pm3"`
	Tax        float64 `json:"tax"`
}

const fName = "data_conf.json"

var conf gpconf
var saveToFileFlag bool = false
var mutex sync.Mutex

//Cargo availabe in ship
func Cargo() float64 {
	return conf.Cargo
}

//WappPhone will return the Whatsapp phone number to message
func WappPhone() string {
	return conf.MPhoneNum
}

//Tax value charged in stations
func Tax() float64 {
	return conf.Tax
}

//MaxInvest is the maximum amount available to invest in a shopping list
func MaxInvest() float64 {
	return conf.MaxInvest
}

//MessageThreshold is the minimum profit required to send a whatsapp message
func MessageThreshold() float64 {
	return conf.MThreshold
}

//Minpm3 Minimal expected profit amount pem cubic meter of cargo
func Minpm3() float64 {
	return conf.Minpm3
}

func backup() bool {
	fmt.Printf("Failed to open %s\n", fName)
	conf = gpconf{100.0, "data_items.json", "data_locations.json", "data_regions.json", 100000000.0, 1000000000, "558387680888", 20, 100000, 0.05}
	return true
}
func init() {
	mutex = sync.Mutex{}
	raw, err := ioutil.ReadFile(fName)
	_ = (err == nil && json.Unmarshal(raw, &conf) == nil) || backup()
	utils.WappMessage(WappPhone(), fmt.Sprintf("eve profit threshold:\n%s", utils.FormatCommas(MessageThreshold())))
}

//Terminate method will save possible changes to the configuration file
func Terminate() {
	saveToFileFlag = saveToFileFlag && utils.Save(fName, conf) && !saveToFileFlag
}
