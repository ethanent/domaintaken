package main

import (
	"bufio"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

var domainTLDLocateRegex = regexp.MustCompile(`^.+\.(.+)$`)

var tlds []string
var tldsMux = &sync.RWMutex{}

func ensureFetchedTLDs() {
	tldsMux.Lock()
	defer tldsMux.Unlock()

	if tlds != nil {
		return
	}

	resp, err := http.Get("https://data.iana.org/TLD/tlds-alpha-by-domain.txt")
	defer resp.Body.Close()

	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		s := scanner.Text()

		tlds = append(tlds, strings.ToLower(s))
	}
}

func validateTLD(tld string) bool {
	ensureFetchedTLDs()

	tldsMux.RLock()
	defer tldsMux.RUnlock()

	for _, t := range tlds {
		if tld == t {
			return true
		}
	}

	return false
}

func validateDomainTLD(d string) bool {
	r := domainTLDLocateRegex.FindStringSubmatch(d)

	if len(r) < 2 {
		panic("invalid domain " + d)
	}

	tld := strings.ToLower(r[1])

	return validateTLD(tld)
}
