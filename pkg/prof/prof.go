package prof

import (
	"encoding/json"
	"log"
	"math"
	"runtime"
	"strings"
	"time"
)

type ProfileData struct {
	Min   int64 `json:"min"`
	Max   int64 `json:"max"`
	Total int64 `json:"total"`
	Count int64 `json:"count"`
}

type ProfChannelData struct {
	key   string
	value int64
}

var (
	record      map[string]*ProfileData
	profChannel chan *ProfChannelData
)

func InitProfile() {
	record = make(map[string]*ProfileData, 100)
	profChannel = make(chan *ProfChannelData, 15)
	go func() {
		for {
			select {
			case t := <-profChannel:
				if _, check := record[t.key]; !check {
					record[t.key] = &ProfileData{
						Min:   math.MaxInt,
						Max:   0,
						Total: 0,
						Count: 0,
					}
				}
				data := record[t.key]
				if t.value > data.Max {
					data.Max = t.value
				}

				if t.value < data.Min {
					data.Min = t.value
				}

				data.Total += t.value
				data.Count++
			}
		}
	}()
}

func write(key string, value time.Duration) {
	profChannel <- &ProfChannelData{key, value.Nanoseconds()}
}

func Write() func() {
	start := time.Now()
	rpc := make([]uintptr, 1)
	n := runtime.Callers(2, rpc)
	if n < 1 {
		return nil
	}
	frame, _ := runtime.CallersFrames(rpc).Next()
	return func() {
		write(frame.Function, time.Now().Sub(start))
	}
}

func Read() {
	jsonData := make(map[string]ProfileData)
	log.Println(record)

	for k, v := range record {
		apiSplit := strings.Split(k, "/")
		apiStr := apiSplit[len(apiSplit)-1]
		jsonData[apiStr] = ProfileData{
			Min:   v.Min / 1000,
			Max:   v.Max / 1000,
			Total: (v.Total / v.Count) / 1000,
			Count: v.Count,
		}
		log.Println(jsonData[apiStr])

	}
	jsonOutput, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		log.Println("Error marshalling to JSON:", err)
		return
	}

	log.Println(string(jsonOutput))
}
