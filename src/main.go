package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

var (
	NUMBEROFCONNECTIONS = 100 // Number of concurrent connections to make to each server
)

func main() {

	addrs := []string{
		"server.com",
	}

	var wg sync.WaitGroup

	for _, addr := range addrs {
		for i := 1; i <= NUMBEROFCONNECTIONS; i++ {
			wg.Add(1)
			go func(addr string, connectionNumber int) {
				defer wg.Done()

				for {
					ipAddrs, err := net.LookupIP(addr)
					if err != nil {
						fmt.Println("Error resolving address", addr, err)
						return
					}

					ipAddr := ipAddrs[0]

					conn, err := net.Dial("tcp", fmt.Sprintf("%s:443", ipAddr))
					if err != nil {
						fmt.Println("Error connecting to", addr, err)
						return
					}

					defer conn.Close()
					fmt.Printf("Connected to %s (Connection #%d)\n", addr, connectionNumber)

					client := &http.Client{}
					req, err := http.NewRequest("GET", fmt.Sprintf("https://%s", addr), nil)
					if err != nil {
						fmt.Println("Error creating request", addr, err)
						return
					}

					resp, err := client.Do(req)
					if err != nil {
						fmt.Println("Error sending GET request to", addr, err)
						return
					}

					defer resp.Body.Close()
					fmt.Printf("GET request sent to %s (Connection #%d)\n", addr, connectionNumber)
				}
			}(addr, i)
		}
	}

	// Make sure all goroutines are running before sleeping
	time.Sleep(10 * time.Second)

	// Wait for all goroutines to finish after 10 seconds
	wg.Wait()
}
