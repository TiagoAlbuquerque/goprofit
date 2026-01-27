package items

import (
	"goprofit/orders"
	"goprofit/utils"
	"strings"
	"sync"

	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Item mimics the structure of an EVE Online ESI item
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

	BuyOrders  []int64    `json:"-"` // Transient, not serialized
	SellOrders []int64    `json:"-"` // Transient, not serialized
	mu         sync.Mutex `json:"-"` // Protects BuyOrders and SellOrders
}

const itemURL = "https://esi.evetech.net/latest/universe/types/%d"
const fName = "data_items.json"

var items map[int]*Item
var saveToFileFlag bool = false
var mutex sync.Mutex // Only protects write operations

func getItemInfo(id int) *Item {
	// println()
	// println("new Item")
	url := fmt.Sprintf(itemURL, id)
	// println(url)
	var item Item
	err := utils.JSONFromURL(url, &item)
	if err != nil {
		fmt.Printf("Error fetching item %d: %v\n", id, err)
	}
	// println(item.Name)
	item.BuyOrders = []int64{}
	item.SellOrders = []int64{}
	// Removed side effects: items[id] = &item and saveToFileFlag = true
	return &item
}

// Get will return the item specified by the provided itemID
func Get(itemID int) *Item {
	// Lock-free read - items data is stable once written
	item, ok := items[itemID]
	if ok {
		return item
	}

	// Item not found, fetch it without lock
	newItem := getItemInfo(itemID)

	// Only lock for write
	mutex.Lock()
	items[itemID] = newItem
	saveToFileFlag = true
	mutex.Unlock()
	return newItem
}

// Search returns items matching the query string
func Search(query string) []*Item {
	// Lock-free read - items map is stable
	query = strings.ToLower(query)
	var results []*Item

	for _, item := range items {
		if strings.Contains(strings.ToLower(item.Name), query) {
			results = append(results, item)
			if len(results) >= 50 {
				break
			}
		}
	}
	return results
}

func (item *Item) place(o orders.Order) {
	item.mu.Lock()
	defer item.mu.Unlock()
	if o.IsBuyOrder {
		item.BuyOrders = append(item.BuyOrders, o.OrderID)
	} else {
		item.SellOrders = append(item.SellOrders, o.OrderID)
	}
}

// IsOfficer will check if the item is an office type item
func (item *Item) IsOfficer() bool {
	for _, v := range item.DogmaAttributes {
		if v.AttributeID == 1692 && v.Value == 5.0 {
			return true
		}
	}
	return false
}

// PlaceOrder will put the received market order in the item list
func PlaceOrder(o orders.Order) {
	item := Get(o.ItemID)
	item.place(o)
}

// Cleanup will clear all the items orders
func Cleanup() {
	// Lock-free iteration - items map is stable
	for _, item := range items {
		item.mu.Lock()
		item.BuyOrders = []int64{}
		item.SellOrders = []int64{}
		item.mu.Unlock()
	}
}

func backup() bool {
	fmt.Printf("Failed to open %s\n", fName)
	items = make(map[int]*Item)
	return true
}

func init() {
	// Try loading from current dir
	raw, err := ioutil.ReadFile(fName)
	if err != nil {
		// Try loading from parent dir
		raw, err = ioutil.ReadFile("../" + fName)
	}

	_ = (err == nil && json.Unmarshal(raw, &items) == nil) || backup()
	Cleanup()
}

// Terminate will save the items files if there are any new items
func Terminate() {
	mutex.Lock()
	defer mutex.Unlock()
	if saveToFileFlag {
		err := utils.SaveToJSONFile(fName, items)
		if err != nil {
			fmt.Println("Failed to save items to file:", err)
		}
		saveToFileFlag = false
	}
}
