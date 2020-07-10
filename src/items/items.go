package items

import (
	"../orders"
	"../utils"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
)

//Item mimics the structure of an EVE Online ESI item
type Item struct {
	Capacity        float32 `json:"capacity"`
	Description     string  `json:"description"`
	DogmaAttributes []struct {
		AttributeID int     `json:"attribute_id"`
		Value       float32 `json:"value"`
	} `json:"dogma_attributes"`
	DogmaEffects []struct {
		EffectID  int  `json:"effect_id"`
		IsDefault bool `json:"is_default"`
	} `json:"dogma_effects"`
	GraphicID      int     `json:"graphic_id"`
	GroupID        int     `json:"group_id"`
	MarketGroupID  int     `json:"market_group_id"`
	Mass           float32 `json:"mass"`
	Name           string  `json:"name"`
	PackagedVolume float32 `json:"packaged_volume"`
	PortionSize    int     `json:"portion_size"`
	Published      bool    `json:"published"`
	Radius         float32 `json:"radius"`
	ItemID         int     `json:"type_id"`
	Volume         float64 `json:"volume"`

	BuyOrders  []int64
	SellOrders []int64
}

const itemURL = "https://esi.evetech.net/latest/universe/types/%d"
const fName = "data_items.json"

var items map[int]*Item
var saveToFileFlag bool = false
var mutex sync.Mutex

func getItemInfo(id int) *Item {
	println()
	println("new Item")
	url := fmt.Sprintf(itemURL, id)
	println(url)
	var item Item
	utils.JSONFromURL(url, &item)
	println(item.Name)
	item.BuyOrders = []int64{}
	item.SellOrders = []int64{}
	items[id] = &item
	saveToFileFlag = true
	return &item
}

//Get will return the item specified by the provided itemID
func Get(itemID int) *Item {
	mutex.Lock()
	defer mutex.Unlock()
	item, ok := items[itemID]
	if !ok {
		item = getItemInfo(itemID)
	}
	return item
}

func (item *Item) place(o orders.Order) {
	if o.IsBuyOrder {
		item.BuyOrders = append(item.BuyOrders, o.OrderID)
	} else {
		item.SellOrders = append(item.SellOrders, o.OrderID)
	}
}

//IsOfficer will check if the item is an office type item
func (item *Item) IsOfficer() bool {
	for _, v := range item.DogmaAttributes {
		if v.AttributeID == 1692 && v.Value == 5.0 {
			return true
		}
	}
	return false
}

//PlaceOrder will put the received market order in the item list
func PlaceOrder(o orders.Order) {
	item := Get(o.ItemID)
	item.place(o)
}

//Cleanup will clear all the items orders
func Cleanup() {
	for _, item := range items {
		item.BuyOrders = []int64{}
		item.SellOrders = []int64{}
	}
}

func backup() bool {
	fmt.Printf("Failed to open %s\n", fName)
	items = make(map[int]*Item)
	return true
}

func init() {
	mutex = sync.Mutex{}
	raw, err := ioutil.ReadFile(fName)
	_ = (err == nil && json.Unmarshal(raw, &items) == nil) || backup()
	Cleanup()
}

//Terminate will save the items files if there are any new items
func Terminate() {
	saveToFileFlag = saveToFileFlag && utils.Save(fName, items) && !saveToFileFlag
}
