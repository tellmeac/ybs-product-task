package types

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"
)

const (
	hourFormat = "15:04"
)

// IntervalsCollide returns true if any of two intervals are collides.
func IntervalsCollide(intervals ...Interval) bool {
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i].From.Before(intervals[j].From)
	})

	for i := 0; i < len(intervals)-1; i++ {
		if intervals[i].To.After(intervals[i+1].From) || intervals[i].To == intervals[i+1].From {
			return true
		}
	}

	return false
}

// Interval представляет часы работы сотрудника службы доставки.
// Quote: График работы задается списком строк формата HH:MM-HH:MM.
type Interval struct {
	From time.Time
	To   time.Time
}

func (interval *Interval) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("%s-%s",
		interval.From.Format(hourFormat),
		interval.To.Format(hourFormat)))
}

func (interval *Interval) UnmarshalJSON(bytes []byte) error {
	var (
		err error
		raw string
	)

	if err := json.Unmarshal(bytes, &raw); err != nil {
		return nil
	}

	fromTo := strings.SplitN(raw, "-", 2)

	interval.From, err = time.Parse(hourFormat, fromTo[0])
	if err != nil {
		return err
	}

	interval.To, err = time.Parse(hourFormat, fromTo[1])
	if err != nil {
		return err
	}

	return nil
}
