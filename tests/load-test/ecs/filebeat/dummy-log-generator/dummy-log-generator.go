package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/pkg/profile"
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
	Files       int
	StepsPerLog int
	LogsPerFile int
	NgRate      float64
	OutputDir   string
	Parallel    int
}

type Result struct {
	Logs    []Log
	Options *Options
}

func SelectNg(totalCount int, ngCount int) map[int]bool {
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
	steps := make([]Step, options.StepsPerLog)
	stepTemplates := make([]StepTemplate, options.StepsPerLog)
	stepsNumbers := make([]int, options.StepsPerLog)
	for i := 0; i < options.StepsPerLog; i++ {
		stepsNumbers[i] = rand.Intn(options.StepsPerLog * 10)
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
	logs := make([]Log, options.LogsPerFile)
	for logIndex := 0; logIndex < len(logs); logIndex++ {
		logs[logIndex].Steps = make([]Step, options.StepsPerLog)
		logs[logIndex].LogTemplate = log.LogTemplate
		copy(logs[logIndex].Steps, log.Steps)
	}
	result := &Result{Logs: logs}
	result.Options = &options
	return result
}

func Clone(result *Result) *Result {
	logs := make([]Log, len(result.Logs))
	for logIndex := 0; logIndex < len(result.Logs); logIndex++ {
		logs[logIndex].Steps = make([]Step, len(result.Logs[logIndex].Steps))
		logs[logIndex].LogTemplate = result.Logs[logIndex].LogTemplate
		copy(logs[logIndex].Steps, result.Logs[logIndex].Steps)
	}
	cloneResult := &Result{Logs: logs}
	cloneResult.Options = result.Options
	return cloneResult
}

func DetectData(step *Step, isNgStep bool) {
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

	ngCount := int(result.Options.NgRate * float64(result.Options.LogsPerFile))
	isNg := SelectNg(result.Options.LogsPerFile, ngCount)
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

func CreateCsv(result *Result, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		fmt.Println("fail to read file")
		os.Exit(1)
	}

	var byteLog = make([]byte, 0, result.Options.LogsPerFile*result.Options.StepsPerLog*100)
	for logIndex := 0; logIndex < len(result.Logs); logIndex++ {

		byteLog = append(byteLog, "Mode,"...)
		byteLog = append(byteLog, result.Logs[logIndex].LogTemplate.Mode...)
		byteLog = append(byteLog, "\nTesterName,"...)
		byteLog = append(byteLog, result.Logs[logIndex].LogTemplate.Name...)
		byteLog = append(byteLog, "\nDate,"...)
		byteLog = append(byteLog, result.Logs[logIndex].Date...)
		byteLog = append(byteLog, "\nResult,"...)
		byteLog = append(byteLog, result.Logs[logIndex].Result...)
		byteLog = append(byteLog, "\nStep,TstName,LoLimit,Data,UpLimit,Unit,Judge\n"...)

		steps := result.Logs[logIndex].Steps
		for stepIndex := 0; stepIndex < len(steps); stepIndex++ {
			step := steps[stepIndex]
			byteLog = append(byteLog, step.StepTemplate.StepNumber...)
			byteLog = append(byteLog, ',')
			byteLog = append(byteLog, step.StepTemplate.TestName...)
			byteLog = append(byteLog, ',')
			byteLog = append(byteLog, step.StepTemplate.LoLimitString...)
			byteLog = append(byteLog, ',')
			byteLog = append(byteLog, step.DataString...)
			byteLog = append(byteLog, ',')
			byteLog = append(byteLog, step.StepTemplate.UpLimitString...)
			byteLog = append(byteLog, ',')
			byteLog = append(byteLog, step.StepTemplate.Unit...)
			byteLog = append(byteLog, ',')
			byteLog = append(byteLog, step.Judge...)
			byteLog = append(byteLog, '\n')
		}
		byteLog = append(byteLog, "END\n"...)
	}
	f.Write(byteLog)
}

func main() {
	defer profile.Start(profile.ProfilePath(".")).Stop()

	var (
		f = flag.Int("f", 1, "Number of generate files.\n0 < f <= 10")
		l = flag.Int("l", 10, "Number of logs per file.\n0 < l , s * l <= 10,000,000")
		s = flag.Int("s", 10, "Number of steps per log.\n0 < s , s * l <= 10,000,000")
		n = flag.Float64("n", 0.1, "NG rate.\n0.0 <= n <= 1.0")
		o = flag.String("o", "./", "Output dir.")
		p = flag.Int("p", 1, "Number of parallel executions.\n0 < p <= cpus")
	)
	flag.Parse()
	if *f <= 0 || 10 < *f {
		flag.Usage()
		return
	}
	if *l <= 0 {
		flag.Usage()
		return
	}
	if *s <= 0 {
		flag.Usage()
		return
	}
	if *s**l > 10000000 {
		flag.Usage()
		return
	}
	if *n < 0.0 || 1.0 < *n {
		flag.Usage()
		return
	}
	cpus := runtime.NumCPU()
	if *p > cpus || *p < 1 {
		fmt.Printf("CPUs : %d\n", cpus)
		fmt.Printf("You need to set # of parallel execution less than or equal to %d.\n\n", cpus)
		flag.Usage()
		return
	}

	options := Options{Files: *f, LogsPerFile: *l, StepsPerLog: *s, NgRate: *n, OutputDir: *o, Parallel: *p}

	rand.Seed(time.Now().UnixNano())

	result := New(options)
	var wg sync.WaitGroup
	semaphore := make(chan int, options.Parallel)
	for i := 0; i < options.Files; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			semaphore <- 1
			cloneResult := Clone(result)
			Generate(cloneResult)
			filename := options.OutputDir + "/log_" + strconv.Itoa(i) + ".csv"
			CreateCsv(cloneResult, filename)
			fmt.Printf("Done : %s\n", filename)
			<-semaphore
		}(i)
	}
	wg.Wait()

	fmt.Println("finish.")
}
