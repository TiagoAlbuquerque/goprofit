package utils

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/goccy/go-json"

	"time"

	"goprofit/utils/color"

	"github.com/cheggaaa/pb/v3"
)

const (
	retryCount    = 5
	retryInterval = time.Second
)

// httpClient is a reusable HTTP client with connection pooling for better performance.
var httpClient *http.Client

func init() {
	transport := &http.Transport{
		MaxIdleConns:        200,
		MaxIdleConnsPerHost: 200,
		IdleConnTimeout:     90 * time.Second,
		ForceAttemptHTTP2:   true,
		DisableCompression:  false, // Explicitly false
	}
	httpClient = &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}
}

func GetURL(url string) (*http.Response, error) {
	var (
		res *http.Response
		err error
	)
	for i := 0; i < retryCount; i++ {
		res, err = httpClient.Get(url)
		if err == nil {
			return res, nil
		}
		log.Printf("error getting URL %s: %v", url, err)
		time.Sleep(retryInterval)
	}
	return nil, fmt.Errorf("failed to get URL %s after %d retries", url, retryCount)
}

func JSONFromURL(url string, result interface{}) error {
	response, err := GetURL(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return json.NewDecoder(response.Body).Decode(result)
}

func writeToFile(filePath string, data []byte) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
}

func marshalToJSONIndent(obj interface{}) ([]byte, error) {
	return json.MarshalIndent(obj, "", "  ")
}

func SaveToJSONFile(filePath string, obj interface{}) error {
	data, err := marshalToJSONIndent(obj)
	if err != nil {
		log.Printf("error marshalling data for %s: %v", filePath, err)
		return err
	}
	err = writeToFile(filePath, data)
	if err != nil {
		log.Printf("error writing to file %s: %v", filePath, err)
		return err
	}
	fmt.Println(filePath, " saved")
	return nil
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

// Top will get the first few items of a sortable list in order
func Top(list sortable) {
	qTop(list, 0, list.Len()-1)
}

// StatusLine will print the providade text in the status indicator line
func StatusLine(c int, text string) {
	up := "\033[A"  //move cursor up one line
	cr := "\r"      //carriage return [volta para o início]
	cl := "\033[2K" //clear line

	fmt.Printf(up+cr+cl+"%s\n", color.Fg8b(c, text))
	//gui.StatusLabel(text)
}

// ProgressBar will start a new progressbar
func ProgressBar(total int, c chan interface{}) {
	tmpl := `{{with string . "prefix"}}{{.}} {{end}}{{counters . }} {{bar . "[" "=" ">" " " "]"}} {{percent . }} {{speed . }} {{rtime . "ETA %s"}}{{with string . "suffix"}} {{.}}{{end}}`

	bar := pb.ProgressBarTemplate(tmpl).Start(total)
	for i := 0; i < total; i++ {
		<-c
		bar.Increment()
	}
	bar.Finish()
}

// KMB will format will compress long values into K M or B formats
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

// FormatCommas will produce spaces at every power of 1000 as 1 000 000 000
func FormatCommas(num float64) string {
	parts := strings.Split(fmt.Sprintf("%.2f", num), ".")
	return commas(parts[0]) + "." + parts[1]
}
