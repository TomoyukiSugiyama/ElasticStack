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
	StepTemplate *StepTemplate
	Data         interface{}
	DataString   string
	Judge        string
}

type StepTemplate struct {
	StepNumber    string
	TestName      string
	LoLimit       interface{}
	LoLimitString string
	UpLimit       interface{}
	UpLimitString string
	Unit          string
}

type Log struct {
	LogTemplate *LogTemplate
	Date        string
	Result      string
	Steps       []Step
}

type LogTemplate struct {
	Mode string
	Name string
}

type Options struct {
	StepCount int
	LogCount  int
	NgRate    float64
}

type Result struct {
	Logs    []Log
	Options *Options
}

func SelectNg(totalCount int, ngCount int) map[int]bool {
	t := time.Now()
	rand.Seed(t.UnixNano())
	isNg := make(map[int]bool)
	for i := 0; i < ngCount; {
		n := rand.Intn(totalCount)
		if !isNg[n] {
			isNg[n] = true
			i++
		}
	}
	return isNg
}

func New(options Options) *Result {
	t := time.Now()
	rand.Seed(t.UnixNano())
	steps := make([]Step, options.StepCount)
	stepTemplates := make([]StepTemplate, options.StepCount)
	stepsNumbers := make([]int, options.StepCount)
	for i := 0; i < options.StepCount; i++ {
		stepsNumbers[i] = rand.Intn(options.StepCount * 10)
	}
	sort.Ints(stepsNumbers)
	units := []string{"V", "MV", "A", "MA", "HEX", "MS"}
	for i := 0; i < len(steps); i++ {
		steps[i].StepTemplate = &stepTemplates[i]
		stepTemplates[i].StepNumber = strconv.Itoa(stepsNumbers[i])
		stepTemplates[i].TestName = "test_" + strconv.Itoa(i)
		unitIndex := rand.Intn(len(units))
		stepTemplates[i].Unit = units[unitIndex]
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
			stepTemplates[i].LoLimit = loLimit
			stepTemplates[i].LoLimitString = fmt.Sprintf("%X", loLimit)
			stepTemplates[i].UpLimit = upLimit
			stepTemplates[i].UpLimitString = fmt.Sprintf("%X", upLimit)
		} else {
			loLimit := rand.Float64() * digitFloat
			upLimit := rand.Float64() * digitFloat
			if loLimit > upLimit {
				tmp := loLimit
				loLimit = upLimit
				upLimit = tmp
			}
			stepTemplates[i].LoLimit = loLimit
			stepTemplates[i].LoLimitString = fmt.Sprintf("%.3f", loLimit)
			stepTemplates[i].UpLimit = upLimit
			stepTemplates[i].UpLimitString = fmt.Sprintf("%.3f", upLimit)
		}
	}

	logTemplate := &LogTemplate{Mode: "dev", Name: "dummy"}
	log := Log{LogTemplate: logTemplate, Steps: steps}
	logs := make([]Log, options.LogCount)
	for logIndex := 0; logIndex < len(logs); logIndex++ {
		logs[logIndex].Steps = make([]Step, options.StepCount)
		logs[logIndex].LogTemplate = log.LogTemplate
		copy(logs[logIndex].Steps, log.Steps)
	}
	result := &Result{Logs: logs}
	result.Options = &options
	return result
}

func DetectData(step *Step, isNgStep bool) {
	t := time.Now()
	rand.Seed(t.UnixNano())
	isHiNg := rand.Intn(2) == 0
	step.Judge = "OK"
	stepTemplate := step.StepTemplate
	if stepTemplate.Unit == "HEX" {
		var data int
		if stepTemplate.UpLimit.(int) == stepTemplate.LoLimit.(int) {
			data = stepTemplate.UpLimit.(int)
		} else {
			data = rand.Intn(stepTemplate.UpLimit.(int)-stepTemplate.LoLimit.(int)) + stepTemplate.LoLimit.(int)
		}
		if isNgStep && isHiNg {
			data = data + stepTemplate.UpLimit.(int)
			step.Judge = "HI"
		}
		if isNgStep && !isHiNg {
			data = data - stepTemplate.UpLimit.(int)
			step.Judge = "LO"
		}
		step.Data = data
		step.DataString = fmt.Sprintf("%X", data)
	} else {
		data := rand.Float64()*(stepTemplate.UpLimit.(float64)-stepTemplate.LoLimit.(float64)) + stepTemplate.LoLimit.(float64)
		if isNgStep && isHiNg {
			data = data + stepTemplate.UpLimit.(float64)
			step.Judge = "HI"
		}
		if isNgStep && !isHiNg {
			data = data - stepTemplate.UpLimit.(float64)
			step.Judge = "LO"
		}
		step.Data = data
		step.DataString = fmt.Sprintf("%.3f", data)
	}
}
func GenerateSteps(steps []Step, isNgLog bool) {
	t := time.Now()
	rand.Seed(t.UnixNano())
	stepsCount := len(steps)
	var ngCount int
	if stepsCount <= 1 {
		ngCount = 1
	} else {
		ngCount = rand.Intn(stepsCount-1) + 1
	}
	isNgStep := SelectNg(stepsCount, ngCount)
	for stepIndex := 0; stepIndex < len(steps); stepIndex++ {
		isNg := isNgLog && isNgStep[stepIndex]
		DetectData(&steps[stepIndex], isNg)
	}
}

func Generate(result *Result) {
	const dayLayout = "2006/01/02,15:04:05"
	t := time.Now()

	ngCount := int(result.Options.NgRate * float64(result.Options.LogCount))
	isNg := SelectNg(result.Options.LogCount, ngCount)
	for logIndex := 0; logIndex < len(result.Logs); logIndex++ {
		if isNg[logIndex] {
			result.Logs[logIndex].Result = "NG"
		} else {
			result.Logs[logIndex].Result = "OK"
		}
		t = t.Add(5 * time.Minute)
		result.Logs[logIndex].Date = t.Format(dayLayout)
		GenerateSteps(result.Logs[logIndex].Steps, isNg[logIndex])
	}
}

func CreateCsv(result *Result) {
	for logIndex := 0; logIndex < len(result.Logs); logIndex++ {
		fmt.Printf("Mode,%s\n", result.Logs[logIndex].LogTemplate.Mode)
		fmt.Printf("TesterName,%s\n", result.Logs[logIndex].LogTemplate.Name)
		fmt.Printf("Date,%s\n", result.Logs[logIndex].Date)
		fmt.Printf("Result,%s\n", result.Logs[logIndex].Result)
		fmt.Printf("Step,TstName,LoLimit,Data,UpLimit,Unit,Judge\n")
		steps := result.Logs[logIndex].Steps
		for stepIndex := 0; stepIndex < len(steps); stepIndex++ {
			fmt.Printf("%s,", steps[stepIndex].StepTemplate.StepNumber)
			fmt.Printf("%s,", steps[stepIndex].StepTemplate.TestName)
			fmt.Printf("%s,", steps[stepIndex].StepTemplate.LoLimitString)
			fmt.Printf("%s,", steps[stepIndex].DataString)
			fmt.Printf("%s,", steps[stepIndex].StepTemplate.UpLimitString)
			fmt.Printf("%s,", steps[stepIndex].StepTemplate.Unit)
			fmt.Printf("%s\n", steps[stepIndex].Judge)
		}
		fmt.Printf("END\n")
	}
}

func main() {
	var (
		s = flag.Int("s", 10, "step count (0 < s)")
		l = flag.Int("l", 10, "log count (0 < l)")
		n = flag.Float64("n", 0.1, "ng rate (0 <= n <= 1)")
	)
	flag.Parse()
	if *s <= 0 {
		flag.Usage()
		return
	}
	if *l <= 1 {
		flag.Usage()
		return
	}
	if *n < 0 || *n > 1 {
		flag.Usage()
		return
	}

	options := Options{StepCount: *s, LogCount: *l, NgRate: *n}
	result := New(options)
	Generate(result)
	CreateCsv(result)
}
