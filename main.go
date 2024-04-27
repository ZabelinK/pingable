package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/go-ping/ping"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	var hosts []string
	for scanner.Scan() {
		hosts = append(hosts, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "reading standard input: %s\n", err)
		return
	}

	results := make(chan string)
	var wg sync.WaitGroup

	// Start a goroutine for each host
	for _, host := range hosts {
		wg.Add(1)
		go func(hostname string) {
			defer wg.Done()
			if reachable, _ := isHostPingable(hostname); reachable {
				results <- hostname
			}
		}(host)
	}

	// Close results channel when all ping operations are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	for hostname := range results {
		fmt.Println(hostname)
	}
}

// isHostPingable uses ICMP Echo Request to determine if the host is pingable.
func isHostPingable(hostname string) (bool, error) {
	pinger, err := ping.NewPinger(hostname)
	if err != nil {
		return false, err
	}

	pinger.Count = 3
	pinger.Timeout = time.Second * 10

	// Run the pinger
	err = pinger.Run()
	if err != nil {
		return false, err
	}

	stats := pinger.Statistics() // Get the statistics
	return stats.PacketsRecv > 0, nil
}
