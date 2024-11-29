package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

// Initialize logging
func StartLogging(fileName string, currFile *os.File) *os.File {
	if currFile != nil {
		log.Println("closing log")
		currFile.Close()
	}

	fh, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return nil
	}

	log.SetOutput(fh)
	if *Debug {
		log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.LUTC | log.Ldate | log.Ldate)
	} else {
		log.SetFlags(log.Lmicroseconds | log.LUTC | log.Ldate | log.Ldate)
	}

	log.Printf("%s v%s\n", tool, version)
	return fh
}

func Version(b bool) {
	if b {
		if commit != "" {
			// go build -ldflags="-X main.commit=$(git rev-parse --short HEAD) -X main.branch=$(git branch | sed 's/.*\* //')"
			fmt.Printf("%s v%s (commit:%s branch:%s)\n", tool, version, commit, branch)
		} else {
			// go build
			fmt.Printf("%s v%s\n", tool, version)
		}

		os.Exit(0)
	}
}

type Job struct {
	Device   string
	Type     string
	Interval int
}

func LoadJobs(fileName string) (*[]Job, *[]Job, error) {
	fh, err := os.Open(fileName)
	if err != nil {
		return nil, nil, err
	}

	defer fh.Close()

	reader := csv.NewReader(fh)
	reader.Comma = '\t'
	reader.Read() // ditch the header
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	var stateJobs []Job
	var reportJobs []Job
	for _, row := range rows {
		//                     0        1        2
		// columns in row => Device  JobType  Interval
		n, _ := strconv.Atoi(row[2])
		job := Job{Device: row[0], Type: row[1], Interval: n}
		switch row[1] {
		case "state":
			stateJobs = append(stateJobs, job)
		case "report":
			reportJobs = append(reportJobs, job)
		default:
			log.Printf("skipping row %d: unknown job type %s for device %s\n", row[2], row[1], row[0])
		}
	}

	log.Printf("loaded %d state jobs\n", len(stateJobs))
	log.Printf("loaded %d report jobs\n", len(reportJobs))
	return &stateJobs, &reportJobs, nil
}

// Different jobs could have different evaluation intervals.  From the total
// list of jobs, we only want the jobs on the same interval
// for 'm', key = interval from config. value = the count of jobs in the interval
func StateTickerMap(jobs *[]Job) map[int]int {
	m := make(map[int]int)
	for _, v := range *jobs {
		m[v.Interval]++
	}

	for k, v := range m {
		log.Printf("identified %d jobs for interval %d\n", v, k)
	}
	return m
}

// group together those jobs where job.Interval matches 'interval' and return that list
func JobsByInterval(interval int, jobs *[]Job) *[]Job {
	var byInterval []Job
	for _, j := range *jobs {
		if j.Interval == interval {
			byInterval = append(byInterval, j)
		}
	}

	return &byInterval
}
