package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"gopkg.in/cheggaaa/pb.v1"
)

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

func JsonFromUrl(url string, out interface{}) {
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

func Load(f_name string) (interface{}, error) {
	raw, err := ioutil.ReadFile(f_name)
	if err != nil {
		return nil, err
	}

	var c interface{}
	json.Unmarshal(raw, &c)
	return c, nil
}

func Save(f_name string, data interface{}) {
	out, _ := json.MarshalIndent(data, "", "  ")
	f, err := os.Create(f_name)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.Write(out)
	fmt.Println(f_name, " saved")
}

type sortable interface {
	Len() int
	Less(i, j int) bool
	Swap(i, j int)
}

func iSort(list sortable, inicio, fim int) {
	i := inicio + 1
	for i <= fim {
		j := i
		for j > inicio && !list.Less(j-1, j) {
			list.Swap(j-1, j)
			j--
		}
		i++
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

func Top(list sortable) {
	qTop(list, 0, list.Len()-1)
}

func StatusLine(text string) {
	fmt.Printf("\033[A\r\033[K%s\n", text)
}

func ProgressBar(total int, c chan bool) {
	bar := pb.StartNew(total)
	for i := 0; i < total; i++ {
		_ = <-c
		bar.Increment()
	}
	bar.Finish()
}

func commas(s string) string {
	if len(s) <= 3 {
		return s
	} else {
		return commas(s[0:len(s)-3]) + " " + s[len(s)-3:]
	}
}

func FormatCommas(num float64) string {
	parts := strings.Split(fmt.Sprintf("%.2f", num), ".")
	if parts[0][0] == '-' {
		return "-" + commas(parts[0][1:]) + "." + parts[1]
	}
	return commas(parts[0]) + "." + parts[1]
}
