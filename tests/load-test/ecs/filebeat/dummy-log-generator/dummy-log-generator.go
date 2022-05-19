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
		fmt.Printf("digit : %d\n", digit)
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
		//d = flag.Int("d", 60, "duration")
		//p = flag.Int("p", 1, "period")
		n = flag.Int("n", 10, "step count")
	)
	flag.Parse()
	options := Options{StepCount: *n}
	log := New(options)
	for i := 0; i < len(log.Steps); i++ {
		fmt.Printf("%#v\n", log.Steps[i])
	}

	// chStop := make(chan int, 1)
	//create_periodic(*p,chStop)

	// time.Sleep(time.Second * time.Duration(*d))
	// chStop <- 0

	// close(chStop)

	// time.Sleep(time.Second * 1)
}

// func create_log() string {
//     var prefix = "Dummy, 1\n"
//     var detail = ""
//     var suffix = "End\n"

//     t := time.Now()
//     rand.Seed(t.UnixNano())

//     test_step := [8]string{"10", "20", "30", "40", "55","70", "100", "120"}
//     vrange := [8]int{150,2000,99,20000,10,50,2000,80}
//     result := "OK"
//     for i, s := range test_step {
//         value1 := rand.Intn(vrange[i])
//         judge := "OK"
//         if (float64(value1) < float64(vrange[i]) * 0.01 || float64(value1) > float64(vrange[i]) * 0.99){
//             judge = "NG"
//             result = "NG"
//         }
//         detail = detail + s + ",Test" + strconv.Itoa(i+1) + "," + strconv.Itoa(value1) + "," + judge + "\n"
//     }

//     const layout = "2006/01/02 15:04:05"
//     prefix = prefix +
//             "Date," + t.Format(layout) +
//             "\nResult," + result +
//             "\nID,123456\nStep,TestName,Value1,Judge\n"

//     return prefix + detail + suffix
// }

// func output_logs(file_path string,logs string) {
//     file, err := os.OpenFile(file_path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600)
//     if err != nil {
//         log.Fatal(err)
//     }
//     defer file.Close()

//     file.Write(([]byte)(logs))
// }

// func create_periodic(period int,stopTimer chan int) {
//     go func() {
//         ticker := time.NewTicker(time.Duration(period) * time.Second)

//     LOOP:
//         for {
//             select {
//             case <-ticker.C:
//                 const layout = "2006-01-02 15:04:05"
//                 fmt.Print("[" + time.Now().Format(layout) + "]")
//                 logs1 := create_log()
//                 output_logs("../node-01/logs/log-01.csv",logs1)
//                 logs2 := create_log()
//                 output_logs("../node-02/logs/log-01.csv",logs2)
//                 fmt.Println("create log.")
//             case <-stopTimer:
//                 ticker.Stop()
//                 break LOOP
//             }
//         }
//     }()
// }