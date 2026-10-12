package pocketsmith

import (
	"fmt"
	"regexp"
	"strconv"
)

// roundUpRe matches the Up Bank round-up suffix appended to Pocketsmith payee strings.
// e.g. "Woolworths, including -1.00 Round-up" or "Grill'd, including $0.42 Round-up"
var roundUpRe = regexp.MustCompile(`\,\sincluding\s*-?\$?([\d.]+)\s*Round-up`)

// IsRoundUp reports whether the payee string contains an Up Bank round-up suffix.
func IsRoundUp(payee string) bool {
	return roundUpRe.MatchString(payee)
}

// ExtractRoundUp strips the round-up suffix from payee and returns the clean
// payee name, the round-up amount, and true. Returns ("", 0, false) if the
// payee does not contain a round-up suffix.
func ExtractRoundUp(payee string) (cleanPayee string, roundUpAmount float64, ok bool) {
	matches := roundUpRe.FindAllStringSubmatch(payee, -1)
	if len(matches) == 0 || len(matches[0]) < 2 {
		return "", 0, false
	}
	amount, err := strconv.ParseFloat(matches[0][1], 64)
	if err != nil {
		return "", 0, false
	}
	return roundUpRe.ReplaceAllString(payee, ""), amount, true
}

// MustExtractRoundUp is like ExtractRoundUp but panics if the payee does not
// contain a round-up suffix.
func MustExtractRoundUp(payee string) (string, float64) {
	clean, amt, ok := ExtractRoundUp(payee)
	if !ok {
		panic(fmt.Sprintf("pocketsmith: no round-up suffix in payee %q", payee))
	}
	return clean, amt
}
