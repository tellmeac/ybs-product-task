package types

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const (
	hourFormat = "15:04"
)

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
