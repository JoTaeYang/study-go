package prof

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"runtime"
	"strings"
	"time"
)

type ProfileData struct {
	min   int64
	max   int64
	total int64
	count int64
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
						min:   math.MaxInt,
						max:   0,
						total: 0,
						count: 0,
					}
				}
				data := record[t.key]
				if t.value > data.max {
					data.max = t.value
				}

				if t.value < data.min {
					data.min = t.value
				}

				data.total += t.value
				data.count++
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
	log.Println(record)

	buf := bytes.Buffer{}
	buf.Grow(2048)
	buf.WriteString("\nAPI\t\t MIN\t MAX\t AVG\tCount\n")
	for k, v := range record {
		apiSplit := strings.Split(k, "/")
		apiStr := apiSplit[len(apiSplit)-1]
		minStr := fmt.Sprintf("  %d ns", v.min/1000)
		maxStr := fmt.Sprintf("%d ns", v.max/1000)
		avg := fmt.Sprintf("%d ns", (v.total/v.count)/1000)
		str := fmt.Sprintf("%s  %s  %s  %s  %d\n", apiStr, minStr, maxStr, avg, v.count)
		buf.WriteString(str)
	}

	log.Println(buf.String())
}
