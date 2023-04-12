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

// Hour представляет час работы сотрудника службы доставки.
// Quote: График работы задается списком строк формата HH:MM-HH:MM.
type Hour struct {
	From time.Time
	To   time.Time
}

func (h *Hour) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		fmt.Sprintf("%s-%s",
			h.From.Format(hourFormat),
			h.To.Format(hourFormat),
		),
	)
}

func (h *Hour) UnmarshalJSON(bytes []byte) error {
	var (
		err error
		raw string
	)

	if err := json.Unmarshal(bytes, &raw); err != nil {
		return nil
	}

	fromTo := strings.SplitN(raw, "-", 2)

	h.From, err = time.Parse(hourFormat, fromTo[0])
	if err != nil {
		return err
	}

	h.To, err = time.Parse(hourFormat, fromTo[1])
	if err != nil {
		return err
	}

	return nil
}
