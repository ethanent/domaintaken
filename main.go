package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"sync"
)

var maxConcurrent = flag.Int("concurrent", 12, "Maximum number of concurrent DNS requests")

func main() {
	flag.Parse()
	initialDomains := flag.Args()

	var domains []string

	for _, d := range initialDomains {
		variants := generateVariants(d)

		domains = append(domains, variants...)
	}

	wg := &sync.WaitGroup{}
	count := 0

	for _, d := range domains {
		wg.Add(1)
		count++

		go checkDomain(d, wg)

		if count >= *maxConcurrent {
			// Note that this will drop our concurrent down to 0 before proceeding to the next maxConcurrent concurrent lookups.
			// Considering this, there is clearly room for optimization.
			wg.Wait()
		}
	}

	wg.Wait()
}

func checkDomain(d string, wg *sync.WaitGroup) {
	validTLD, err := validateDomainTLD(d)

	if err != nil {
		printColor(color.New(color.BgYellow), d, "ERR", err.Error())

		wg.Done()
		return
	}

	if !validTLD {
		printColor(color.New(color.BgYellow), d, "INVALID TLD")

		wg.Done()
		return
	}

	exist, err := domainExists(d)

	if err != nil {
		printColor(color.New(color.BgYellow), d, "ERR", err.Error())
	} else {
		if exist {
			printColor(color.New(color.BgRed), d, "TAKEN")
		} else {
			printColor(color.New(color.BgGreen), d, "AVAIL")
		}
	}

	wg.Done()
}

var stdoutMux = &sync.Mutex{}

func printColor(c *color.Color, text ...interface{}) {
	stdoutMux.Lock()
	defer stdoutMux.Unlock()

	c.Print(" ")
	fmt.Print(" ")
	fmt.Println(text...)
}
