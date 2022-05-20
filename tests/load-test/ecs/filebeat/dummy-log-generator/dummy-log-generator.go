package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"time"
)

type Step struct {
	StepId   string
	TestName string
	Data     interface{}
	LoLimit  interface{}
	UpLimit  interface{}
	Unit     string
	Judge    string
}

type Log struct {
	Mode   string
	Name   string
	Date   string
	Result string
	Steps  []Step
	Suffix string
}

type Options struct {
	StepCount int
}

func New(options Options) *Log {
	day := time.Now()
	const dayLayout = "2006/01/02,15:04:05"
	date := day.Format(dayLayout)

	rand.Seed(day.UnixNano())
	steps := make([]Step, options.StepCount)
	stepsIds := make([]int, options.StepCount)
	for i := 0; i < options.StepCount; i++ {
		stepsIds[i] = rand.Intn(options.StepCount * 10)
	}
	sort.Ints(stepsIds)
	units := []string{"V", "MV", "A", "MA", "HEX", "MS"}
	for i := 0; i < options.StepCount; i++ {
		steps[i].StepId = strconv.Itoa(stepsIds[i])
		steps[i].TestName = "test_" + strconv.Itoa(i)
		unitIndex := rand.Intn(len(units))
		steps[i].Unit = units[unitIndex]
		digit := rand.Intn(5)
		digitFloat := math.Pow10(digit)

		if units[unitIndex] == "HEX" {
			loLimit := rand.Intn(1 << (digit * 4))
			upLimit := rand.Intn(1 << (digit * 4))

			if loLimit > upLimit {
				tmp := loLimit
				loLimit = upLimit
				upLimit = tmp
			}
			steps[i].LoLimit = fmt.Sprintf("%x", loLimit)
			steps[i].UpLimit = fmt.Sprintf("%x", upLimit)
		} else {
			loLimit := rand.Float64() * digitFloat
			upLimit := rand.Float64() * digitFloat
			if loLimit > upLimit {
				tmp := loLimit
				loLimit = upLimit
				upLimit = tmp
			}
			steps[i].LoLimit = fmt.Sprintf("%.3f", loLimit)
			steps[i].UpLimit = fmt.Sprintf("%.3f", upLimit)
		}
	}
	log := &Log{Mode: "dev", Name: "dummy", Date: date, Steps: steps}
	return log
}

func main() {
	var (
		n = flag.Int("n", 10, "step count")
	)
	flag.Parse()
	options := Options{StepCount: *n}
	log := New(options)
	for i := 0; i < len(log.Steps); i++ {
		fmt.Printf("%#v\n", log.Steps[i])
	}

}
