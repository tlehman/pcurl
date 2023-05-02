package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
)

func downloadFile(url string, concurrency int) error {
	// Get the file size
	resp, err := http.Head(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fileSize := resp.ContentLength
	if fileSize <= 0 {
		return fmt.Errorf("unable to determine file size")
	}

	// Create the output file
	outputFile, err := os.Create("output_file")
	if err != nil {
		return err
	}
	defer outputFile.Close()

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
			defer wg.Done()

			// Create a new HTTP request with the range header
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				fmt.Printf("Error creating request: %v\n", err)
				return
			}
			req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))

			// Execute the request
			resp, err := http.DefaultClient.Do(req)
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

	err = downloadFile(url, concurrency)
	if err != nil {
		fmt.Printf("Error downloading file: %v\n", err)
		return
	}

	fmt.Println("File downloaded successfully!")
}
