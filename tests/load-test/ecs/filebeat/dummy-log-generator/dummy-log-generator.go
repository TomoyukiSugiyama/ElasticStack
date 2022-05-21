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

	LoLimit interface{}
	UpLimit interface{}
	Unit    string
}

type Log struct {
	Date  string
	Data  []interface{}
	Judge string
}

type Template struct {
	Mode   string
	Name   string
	Result string
	Steps  []Step
	Suffix string
}

type Result struct {
	Logs     []Log
	Template *Template
}

type Options struct {
	StepCount int
	LogCount  int
	NgRate    float64
}

func New(options Options) *Result {
	day := time.Now()
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
	template := &Template{Mode: "dev", Name: "dummy", Steps: steps}
	result := &Result{Template: template}
	result.Logs = make([]Log, options.LogCount)
	return result
}
func Generate(options Options, result *Result) {
	const dayLayout = "2006/01/02,15:04:05"
	for logIndex := 0; logIndex < options.LogCount; logIndex++ {
		day := time.Now()
		date := day.Format(dayLayout)
		result.Logs[logIndex].Date = date
		for stepIndex := 0; stepIndex < options.StepCount; stepIndex++ {

		}
	}
}
func main() {
	var (
		s = flag.Int("s", 10, "step count")
		l = flag.Int("l", 10, "log count")
		n = flag.Float64("n", 0.01, "ng rate")
	)
	flag.Parse()
	options := Options{StepCount: *s, LogCount: *l, NgRate: *n}
	result := New(options)
	Generate(options, result)
	for logIndex := 0; logIndex < options.LogCount; logIndex++ {
		fmt.Printf("%#v\n", result.Logs[logIndex].Date)
	}
	for i := 0; i < len(result.Template.Steps); i++ {
		fmt.Printf("%#v\n", result.Template.Steps[i])
	}

}
