package main

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"sync"
)

func main() {
	domains := os.Args[1:]

	wg := &sync.WaitGroup{}

	for _, d := range domains {
		wg.Add(1)

		go checkDomain(d, wg)
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
		printColor(color.New(color.BgYellow), "ERR", err.Error())
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
