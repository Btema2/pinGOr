package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var client = http.Client{
	Timeout: 10 * time.Second,
}

type Result struct {
	URL        string
	Status     string
	Duration   time.Duration
	Err        error
	StatusCode int
}

func (result Result) String() string {
	if result.Err == nil {
		return fmt.Sprintf("URL: %s | Status: %s | in: %v", result.URL, result.Status, result.Duration.Round(time.Millisecond))
	}
	return fmt.Sprintf("URL: %s | Status: %v | in: %v", result.URL, result.Err, result.Duration.Round(time.Millisecond))
}

func checkURL(url string) Result {
	start := time.Now()
	resp, err := client.Get(url)
	if err != nil {
		return Result{
			URL:      url,
			Status:   "DOWN",
			Duration: time.Since(start),
			Err:      err,
		}
	}
	defer resp.Body.Close()

	return Result{
		URL:        url,
		Status:     resp.Status,
		Duration:   time.Since(start),
		StatusCode: resp.StatusCode,
	}
}

func readFile(filename string) ([]string, error) {
	var targets []string
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		target := strings.TrimSpace(scanner.Text())
		if target == "" || strings.HasPrefix(target, "#") {
			continue
		}
		targets = append(targets, target)
	}
	err = scanner.Err()

	return targets, err
}

func main() {
	start := time.Now()
	targets, err := readFile("targets.txt")
	if err != nil {
		log.Fatal(err)
	}

	downCount := 0
	errorCount := 0
	okCount := 0
	results := make(chan Result, len(targets))
	for _, url := range targets {
		go func() {
			results <- checkURL(url)
		}()
	}
	for range targets {
		res := <-results
		fmt.Println(res)
		if res.Err != nil {
			downCount++
		} else if res.StatusCode >= 400 {
			errorCount++
		} else {
			okCount++
		}

	}
	fmt.Println("Finished in", time.Since(start), "| OK:", okCount, "| Error:", errorCount, "| DOWN:", downCount)
}
