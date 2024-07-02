package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

func main() {
	var workers, cases int
	var profiling bool
	var filterRegex string
	flag.IntVar(&workers, "workers", runtime.NumCPU(), "number of workers")
	flag.IntVar(&cases, "cases", 10000, "number of test cases")
	flag.BoolVar(&profiling, "profiling", false, "enable profiling on port 8080")
	flag.StringVar(&filterRegex, "filter", "", "filter test case filename by regex")
	flag.Parse()
	var filterRx *regexp.Regexp
	if filterRegex != "" {
		var err error
		filterRx, err = regexp.Compile(filterRegex)
		if err != nil {
			log.Fatal(err)
		}
	}
	m, _ := mem.VirtualMemory()
	cp, _ := cpu.Info()
	fmt.Printf("Running on %s with %d cores and %dgb mem\n", cp[0].ModelName, runtime.NumCPU(),
		m.Total/uint64(1000000000))
	fmt.Printf("Workers: %d\nCount: %d\n", workers, cases)
	csv := [][]string{
		{"Engine", "Case", "p1", "p2", "p3", "p4", "p5", "Overall"},
	}
	var wg sync.WaitGroup
	if profiling {
		wg.Add(1)
		go func() {
			if err := http.ListenAndServe(":8080", nil); err != nil {
				wg.Done()
				panic(err)
			}
		}()
	}
	for engine, engineIface := range wafInterfaces {
		fmt.Printf("Running tests for %q engine\n", engine)
		waf := engineIface
		engineIface.Init()
		files := []string{
			"../coraza-waf/coraza.conf-recommended",
			"../../../../../projects/coreruleset/crs-setup.conf.example",
			"../../../../../projects/coreruleset/rules/*.conf",
		}
		for _, f := range files {
			// read the file as a glob
			g, err := filepath.Glob(f)
			if err != nil {
				panic(err)
			}
			for _, file := range g {
				if err := waf.LoadDirectives(file); err != nil {
					panic(err)
				}
			}
		}
		testfiles, err := os.ReadDir("./tests")
		if err != nil {
			log.Fatal(err)
		}

		for _, tc := range testfiles {
			if filterRx != nil && !filterRx.MatchString(tc.Name()) {
				continue
			}
			fmt.Println("Opening ", tc.Name())
			c, err := openTest(path.Join("./tests", tc.Name()))
			if err != nil {
				panic(err)
			}
			fmt.Printf("Preparing test case %q on %q\n", c.Name, engine)
			timeStart := time.Now().UnixNano()
			wg := new(sync.WaitGroup)
			var dph [5]int64
			for i := 0; i < workers; i++ {
				wg.Add(1)
				go func(c *testFile) {
					data, err := c.Run(waf, cases)
					if err != nil {
						panic(err)
					}
					for i := range data {
						dph[i] += data[i]
					}
					wg.Done()
				}(c)
			}
			wg.Wait()
			timeEnd := time.Now().UnixNano()
			timeTaken := timeEnd - timeStart
			timePerRequest := timeTaken / int64(cases*workers)
			fmt.Printf("case %q:\n", c.Name)
			fmt.Printf("Took %d seconds, %dus per request.\n", timeTaken/1e9, timePerRequest/1000)
			for i := range dph {
				// now we report the avg duration per phase
				dph[i] /= int64(cases * workers)
				fmt.Printf("Phase %d: %.3fus\n", i+1, float64(dph[i])/1000)
			}
			csv = append(csv, []string{
				engine,
				c.Name,
				fmt.Sprintf("%.3f", float64(dph[0])/1000),
				fmt.Sprintf("%.3f", float64(dph[1])/1000),
				fmt.Sprintf("%.3f", float64(dph[2])/1000),
				fmt.Sprintf("%.3f", float64(dph[3])/1000),
				fmt.Sprintf("%.3f", float64(dph[4])/1000),
				fmt.Sprintf("%.3f", float64(timeTaken)/1000),
			})
			fmt.Println("-----")
		}
	}
	fmt.Println("Writing CSV")
	for _, row := range csv {
		for i, col := range row {
			if i > 0 {
				fmt.Print(",")
			}
			fmt.Print(col)
		}
		fmt.Println()
	}
	if profiling {
		fmt.Println("------")
		fmt.Println("Listening for profiling...")
		wg.Wait()
	}
}
