package parser

import (
	"encoding/csv"
	"fmt"
	"github.com/jedib0t/go-pretty/table"
	"io"
	"math"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type MethodFrequencyStruct struct {
	Method string
	Freq   int
}

type TopHits struct {
	Method string
	Freq   int
	Url    string
}

type LatencyMediator struct {
	Method string
	Times  []float64
}

type Latency struct {
	Method  string
	Url     string
	MinTime float64
	MaxTime float64
	AvgTime float64
}

func LogUploadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("received request to upload csv")

	reader := csv.NewReader(r.Body)
	var items [][]string
	// parse item in the csv string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("error occurred: ", err.Error())
		}

		items = append(items, record)
	}

	fmt.Println("length of csv: ", len(items))

	// Remove first item from the slice as its the header
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	if len(items) > 0 {
		_, result := items[0], items[1:]
		topHits := GetTopURL(result)
		t.AppendHeader(table.Row{"Method", "Url", "Frequency"})
		for _, d := range topHits {
			t.AppendRows([]table.Row{
				{d.Url, d.Freq, d.Method},
			})
		}
		t.Render()

		t = table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Method", "Url", "Min Time", "Max Time", "Avg Time"})
		latencyMetric := GetLatencyMetric(result)
		for _, d := range latencyMetric {
			t.AppendRows([]table.Row{
				{d.Method, d.Url, d.MinTime, d.MaxTime, d.AvgTime},
			})
		}
		t.Render()
	}

	return
}

type kv struct {
	Key   string
	Value int
}

type entries []kv

func (s entries) Len() int           { return len(s) }
func (s entries) Less(i, j int) bool { return s[i].Value < s[j].Value }
func (s entries) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func GetTopURL(data [][]string) []TopHits {
	// Regex to find number
	re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
	topRequestMap := make(map[string]MethodFrequencyStruct, 0)
	freqMap := make(map[string]int)
	result := make([]TopHits, 0)
	var es entries

	for _, i := range data {
		url := i[1] + "_" + i[2]
		if len(url) > 0 {
			if re.Match([]byte(url)) {
				genericVal := re.ReplaceAllString(url, "{id}")
				if len(genericVal) > 0 {
					freqMap[genericVal]++
					if val, ok := topRequestMap[genericVal]; ok {
						val.Freq = val.Freq + 1
						topRequestMap[genericVal] = val
					} else {
						topRequestMap[genericVal] = MethodFrequencyStruct{
							Method: i[2],
							Freq:   1,
						}
					}
				}
			} else {
				freqMap[url]++
				if val, ok := topRequestMap[url]; ok {
					val.Freq = val.Freq + 1
					topRequestMap[url] = val
				} else {
					topRequestMap[url] = MethodFrequencyStruct{
						Method: i[2],
						Freq:   1,
					}
				}
			}
		}
	}

	for k, v := range freqMap {
		es = append(es, kv{Value: v, Key: k})
	}
	sort.Sort(sort.Reverse(es))

	for _, e := range es {
		urlVal := strings.Split(e.Key, "_")
		val := TopHits{
			Method: topRequestMap[e.Key].Method,
			Freq:   e.Value,
			Url:    urlVal[0],
		}
		if len(result) < 5 {
			result = append(result, val)
		} else {
			break
		}
	}

	return result
}

func GetLatencyMetric(data [][]string) []Latency {
	metrics := make(map[string]LatencyMediator)
	result := make([]Latency, 0)
	re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)

	for _, i := range data {
		url := i[1] + "_" + i[2]
		if re.Match([]byte(url)) {
			genericVal := re.ReplaceAllString(url, "{id}")
			if len(genericVal) > 0 {
				if val, ok := metrics[genericVal]; ok {
					flval, err := strconv.ParseFloat(i[3], 64)
					if err != nil {
						fmt.Println("failed to convert string to float 64 ", err.Error())
					}
					val.Times = append(val.Times, flval)
					metrics[genericVal] = val
				} else {
					flval, err := strconv.ParseFloat(i[3], 64)
					if err != nil {
						fmt.Println("failed to convert string to float 64 ", err.Error())
					}
					metrics[genericVal] = LatencyMediator{
						Method: i[2],
						Times:  []float64{flval},
					}
				}
			}
		} else {
			if val, ok := metrics[url]; ok {
				flval, err := strconv.ParseFloat(i[3], 64)
				if err != nil {
					fmt.Println("failed to convert string to float 64 ", err.Error())
				}
				val.Times = append(val.Times, flval)
				metrics[url] = val
			} else {
				flval, err := strconv.ParseFloat(i[3], 64)
				if err != nil {
					fmt.Println("failed to convert string to float 64 ", err.Error())
				}
				metrics[url] = LatencyMediator{
					Method: i[2],
					Times:  []float64{flval},
				}
			}
		}

	}

	for k := range metrics {
		sortedSlice := metrics[k].Times
		lenOfSlice := len(sortedSlice)
		if lenOfSlice > 0 {
			sort.Float64s(sortedSlice)
			min := sortedSlice[0]
			max := sortedSlice[lenOfSlice-1]
			sum := 0.0
			for _, items := range sortedSlice {
				sum = sum + items
			}
			avg := sum / float64(lenOfSlice)

			urlVal := strings.Split(k, "_")
			lncy := Latency{
				Method:  metrics[k].Method,
				Url:     urlVal[0],
				MinTime: min,
				MaxTime: max,
				AvgTime: math.Floor(avg*100)/100,
			}
			result = append(result, lncy)
		}
	}

	return result
}
