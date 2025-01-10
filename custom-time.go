package pocketsmith

import "time"

var customTimeFormat = "2006-01-02"

// customTime is used to unmarshal & marshal dates as the default dates returned
// from the Pocketsmith API are not `ISO 8601` format.
type customTime struct {
	time.Time
}

func (ct *customTime) UnmarshalJSON(b []byte) error {
	str := string(b)
	str = str[1 : len(str)-1] // remove quotes around the date string.
	parsedTime, err := time.Parse(customTimeFormat, str)
	if err != nil {
		return err
	}
	ct.Time = parsedTime
	return nil
}
