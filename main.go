package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

func createOutputFileWithSize(filename string, size int64) error {
	outputFile, err := os.Create(filename)
	defer outputFile.Close()
	if err != nil {
		return err
	}
	for i := 0; i < int(size); i++ {
		outputFile.Write([]byte("0"))
	}
	return nil
}

func getRemoteFileSize(url string) (int64, error) {
	// Get the file size
	resp, err := http.DefaultClient.Head(url)
	if err != nil {
		return int64(-1), err
	}
	defer resp.Body.Close()

	fileSize := resp.ContentLength
	if fileSize <= 0 {
		return int64(-1), fmt.Errorf("unable to determine file size")
	}
	return fileSize, nil
}

func downloadFile(url string, concurrency int) error {

	// Get the file size
	fileSize, err := getRemoteFileSize(url)
	if err != nil {
		return err
	}

	// Calculate the chunk size for each goroutine
	chunkSize := fileSize / int64(concurrency)

	// Create a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup
	wg.Add(concurrency)

	// Launch multiple goroutines to download file in parallel
	for i := 0; i < concurrency; i++ {
		start := int64(i) * chunkSize
		end := start + chunkSize - 1

		// The last goroutine should download the remaining bytes
		if i == concurrency-1 {
			end = fileSize - 1
		}

		go func(start, end int64) {
			transport := http.DefaultTransport.(*http.Transport).Clone()
			transport.DisableKeepAlives = true // make a new connection for each request

			client := &http.Client{
				Transport: transport,
			}
			defer wg.Done()

			// open a new file descriptor for each goroutine and seek to the correct start position
			outputFile, _ := os.OpenFile("output_file", os.O_WRONLY, 0644)
			outputFile.Seek(start, 0)

			// Create a new HTTP request with the range header
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				fmt.Printf("Error creating request: %v\n", err)
				return
			}
			req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))

			// Execute the request
			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("Error executing request: %v\n", err)
				return
			}
			defer resp.Body.Close()

			// Copy the response body to the output file
			_, err = io.Copy(outputFile, resp.Body)
			if err != nil {
				fmt.Printf("Error copying response: %v\n", err)
				return
			}
		}(start, end)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	return nil
}

func main() {
	// Get command-line arguments
	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Println("Usage: pcurl <url> <concurrency>")
		return
	}

	url := args[0]
	concurrency, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println("Invalid concurrency value")
		return
	}

	timeStart := time.Now()
	err = downloadFile(url, concurrency)
	if err != nil {
		fmt.Printf("Error downloading file: %v\n", err)
		return
	}
	timeEnd := time.Now()

	// measure file size
	fileInfo, _ := os.Stat("output_file")
	fileInfo.Size()

	speedMiBps := float64(fileInfo.Size()) / timeEnd.Sub(timeStart).Seconds() / 1024 / 1024

	fmt.Printf(
		"File downloaded in %.0f seconds at %.0f MiB/s!\n",
		timeEnd.Sub(timeStart).Seconds(),
		speedMiBps,
	)
}
