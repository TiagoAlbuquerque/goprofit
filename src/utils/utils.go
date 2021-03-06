package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"../utils/color"

	"../gui"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/Rhymen/go-whatsapp"
	"gopkg.in/cheggaaa/pb.v1"
)

//GetURL will return the httpresponse of a http Get of the provided URL
func GetURL(url string) *http.Response {
	var res *http.Response
	var err error
	for ok := false; !ok; {
		res, err = http.Get(url)
		ok = (err == nil)
		err = nil
	}
	return res
}

//JSONFromURL will unmarshall an JSON received by the http response of the provided URL
func JSONFromURL(url string, out interface{}) {
	var body []byte
	var err error
	for ok := false; !ok; {
		res := GetURL(url)
		defer res.Body.Close()
		body, err = ioutil.ReadAll(res.Body)
		ok = (err == nil)
		if !ok {
			fmt.Print("!OK")
		}
		err = nil
	}
	json.Unmarshal(body, out)
}

//Load will read and unmarshal from JSON the providade file path
func Load(fName string) (interface{}, error) {
	raw, err := ioutil.ReadFile(fName)
	if err != nil {
		return nil, err
	}

	var c interface{}
	json.Unmarshal(raw, &c)
	return c, nil
}

//Save will marshal to JSON and save the providade data to the specified file
func Save(fName string, data interface{}) bool {
	out, _ := json.MarshalIndent(data, "", "  ")
	f, _ := os.Create(fName)
	defer f.Close()
	f.Write(out)
	fmt.Println(fName, " saved")
	return true
}

type sortable interface {
	Len() int
	Less(i, j int) bool
	Swap(i, j int)
}

func iSort(list sortable, inicio, fim int) {
	for i := inicio + 1; i <= fim; i++ {
		for j := i; j > inicio && !list.Less(j-1, j); j-- {
			list.Swap(j-1, j)
		}
	}
}

func qTop(list sortable, inicio, fim int) {
	if fim-inicio <= 25 {
		iSort(list, inicio, fim)
		return
	}
	i := inicio
	j := fim
	p := int((inicio + fim) / 2)
	for i <= j {
		for list.Less(i, p) {
			i++
		}
		for list.Less(p, j) {
			j--
		}
		if i <= j {
			list.Swap(i, j)
			if p == i {
				p = j
			} else if p == j {
				p = i
			}
			i++
			j--
		}
	}
	if inicio < j {
		qTop(list, inicio, j)
	}
}

//Top will get the first few items of a sortable list in order
func Top(list sortable) {
	qTop(list, 0, list.Len()-1)
}

//StatusLine will print the providade text in the status indicator line
func StatusLine(c int, text string) {
	up := "\033[A"  //move cursor up one line
	cr := "\r"      //carriage return [volta para o início]
	cl := "\033[2K" //clear line

	fmt.Printf(up+cr+cl+"%s\n", color.Fg8b(c, text))
	gui.StatusLabel(text)
}

//ProgressBar will start a new progressbar
func ProgressBar(total int, c chan bool) {
	cOK := make(chan bool)
	go gui.ProgressBar(total, cOK)
	bar := pb.StartNew(total)

	for i := 0; i < total; i++ {
		cOK <- <-c
		bar.Increment()
	}
	bar.Finish()
}

//KMB will format will compress long values into K M or B formats
func KMB(num float64) string {
	kmb := ""
	if num >= 1000 {
		num /= 1000
		kmb = "K"
	}
	if num >= 1000 {
		num /= 1000
		kmb = "M"
	}
	if num >= 1000 {
		num /= 1000
		kmb = "B"
	}
	return fmt.Sprintf("%.3f %s", num, kmb)
}

func commas(s string) string {
	if len(s) <= 3 {
		return s
	}
	return commas(s[0:len(s)-3]) + " " + s[len(s)-3:]
}

//FormatCommas will produce spaces at every power of 1000 as 1 000 000 000
func FormatCommas(num float64) string {
	parts := strings.Split(fmt.Sprintf("%.2f", num), ".")
	return commas(parts[0]) + "." + parts[1]
}

var wac, err = whatsapp.NewConn(72 * time.Hour)
var sess whatsapp.Session

func wappInit() {
	qrChan := make(chan string)
	obj := qrcodeTerminal.New()
	go func() {
		obj.Get(<-qrChan).Print()
	}()
	sess, err = wac.Login(qrChan)
}

func init() {
	//wappInit()
}

//WappMessage will send an Whatsapp message
func WappMessage(number, txt string) {
	/*sess, err = wac.RestoreWithSession(sess)
	wac.Send(whatsapp.TextMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: number + "@s.whatsapp.net",
		},
		Text: txt,
	}) //*/
}
