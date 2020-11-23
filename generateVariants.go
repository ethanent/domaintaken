package main

import (
	"regexp"
	"strconv"
)

var replaceRegexp = regexp.MustCompile(`\(([^,)]+)(?:,([^),]+))?\)`)
var alphabet = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
var digits = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

func generateVariants(d string) []string {
	dr := []rune(d)

	submatch := replaceRegexp.FindStringSubmatch(d)

	if submatch == nil || len(submatch) < 2 {
		return []string{d}
	}

	var insertValues []string

	switch submatch[1] {
	case "alpha":
		insertValues = append(insertValues, alphabet...)
	case "num":
		insertValues = append(insertValues, digits...)
	case "alphanum":
		insertValues = append(insertValues, append(alphabet, digits...)...)
	case "tld":
		ensureFetchedTLDs()

		tldsMux.RLock()
		defer tldsMux.RUnlock()

		for _, t := range tlds {
			if submatch[2] != "" {
				preferredLen, err := strconv.ParseInt(submatch[2], 10, 64)

				if err != nil {
					panic(err)
				}

				if len(t) != int(preferredLen) {
					continue
				}
			}

			insertValues = append(insertValues, t)
		}
	default:
		panic("unknown variation method " + submatch[1])
	}

	matchIdx := replaceRegexp.FindStringIndex(d)

	var retVariants []string

	for _, v := range insertValues {
		curVariant := string(dr[0:matchIdx[0]]) + v + string(dr[matchIdx[1]:])
		//retVariants = append(retVariants, curVariant)

		curVariantVariants := generateVariants(curVariant)

		retVariants = append(retVariants, curVariantVariants...)
	}

	return retVariants
}
