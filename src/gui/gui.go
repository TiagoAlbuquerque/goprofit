package gui

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var statusLabel *walk.Label
var pb *walk.ProgressBar
var qrView *walk.ImageView
var tabs *walk.TabWidget

var shopTab = TabPage{Layout: HBox{}, Title: "Shop"}

func init() {
	go MainWindow{
		Title:   "GoProfit",
		MinSize: Size{600, 400},
		Layout:  VBox{},
		MenuItems: []MenuItem{
			Menu{Text: "Items"},
			Menu{Text: "Regions"},
			Menu{Text: "WhatsApp"},
		},
		Children: []Widget{
			TabWidget{
				AssignTo: &tabs,
				Pages: []TabPage{
					shopTab,
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					ProgressBar{AssignTo: &pb},
					Label{AssignTo: &statusLabel, Text: "status", TextAlignment: AlignCenter, MinSize: Size{120, 15}},
				},
			},
		},
	}.Run()
}

func ProgressBar(total int, c chan bool) {
	pb.SetRange(0, total)
	pb.SetValue(0)
	for i := 0; i < total; i++ {
		_ = <-c
		pb.SetValue(i)
	}
}

func StatusLabel(text string) {
	statusLabel.SetText(text)
}
func wAppQRCode() {

}
