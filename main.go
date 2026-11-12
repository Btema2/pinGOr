package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"log"
)

var client = http.Client{
	Timeout: 10 * time.Second,
}

func checkURL(url string) string {
	start := time.Now()
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Sprintf("URL: %s | Status: DOWN | Error: %v | in: %v", url, err, time.Since(start))
	}
	defer resp.Body.Close()

	return fmt.Sprintf("URL: %s | Status: %s | in: %v", url, resp.Status, time.Since(start))
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
		if target == "" || strings.HasPrefix(target, "#") {continue}
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
	for _, url := range targets {
		fmt.Println(checkURL(url))
	}
	fmt.Println("Finished in", time.Since(start))
}
