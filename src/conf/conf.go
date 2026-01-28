package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"

	"goprofit/utils"
)

type gpconf struct {
	Cargo      float64 `json:"cargo"`
	MaxInvest  float64 `json:"max_invest"`
	MThreshold float64 `json:"message_threshold"`
	MPhoneNum  string  `json:"wapp_phone"`
	Minpm3     float64 `json:"min_pm3"`
	Tax        float64 `json:"tax"`
	Threads    int     `json:"threads"` // New field
}

const fName = "data_conf.json"

var conf gpconf
var saveToFileFlag bool = false
var mutex sync.Mutex

// Threads returns the configured number of threads/workers
func Threads() int {
	if conf.Threads <= 0 {
		return 16 // Default safe fallback
	}
	return conf.Threads
}

// Cargo availabe in ship
func Cargo() float64 {
	return conf.Cargo
}

// WappPhone will return the Whatsapp phone number to message
func WappPhone() string {
	return conf.MPhoneNum
}

// Tax value charged in stations
func Tax() float64 {
	return conf.Tax
}

// MaxInvest is the maximum amount available to invest in a shopping list
func MaxInvest() float64 {
	return conf.MaxInvest
}

// MessageThreshold is the minimum profit required to send a whatsapp message
func MessageThreshold() float64 {
	return conf.MThreshold
}

// Minpm3 Minimal expected profit amount pem cubic meter of cargo
func Minpm3() float64 {
	return conf.Minpm3
}

func backup() bool {
	fmt.Printf("Failed to open %s\n", fName)
	conf = gpconf{100.0, 100000000.0, 1000000000, "558387680888", 100000, 0.05, 16}
	return true
}
func init() {
	mutex = sync.Mutex{}
	raw, err := ioutil.ReadFile(fName)
	_ = (err == nil && json.Unmarshal(raw, &conf) == nil) || backup()
}

// Terminate method will save possible changes to the configuration file
func Terminate() {
	if saveToFileFlag {
		err := utils.SaveToJSONFile(fName, conf)
		if err != nil {
			fmt.Println("Failed to save configuration to file:", err)
		}
		saveToFileFlag = false
	}
}
