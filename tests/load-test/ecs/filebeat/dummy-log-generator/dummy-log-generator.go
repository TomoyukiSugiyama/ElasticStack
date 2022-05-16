package main

import (
	"flag"
	"time"
)

type Step struct {
	StepId string
	Value  interface{}
	StdMax interface{}
	StdMin interface{}
	Result string
}

type LogFormat struct {
	Prefix string
	Steps  []Step
	Suffix string
	Result string
}

func main() {
	var (
		d = flag.Int("d", 60, "duration")
		//p = flag.Int("p", 1, "period")
	)
	flag.Parse()

	chStop := make(chan int, 1)
	//create_periodic(*p,chStop)

	time.Sleep(time.Second * time.Duration(*d))
	chStop <- 0

	close(chStop)

	time.Sleep(time.Second * 1)
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
