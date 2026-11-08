package main

import (
	"fmt"
	"net/http"
	"time"
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

func main() {
	fmt.Println(checkURL("https://google.com"))
	fmt.Println(checkURL("https://ejrmjikemrcjikemije.com"))
}