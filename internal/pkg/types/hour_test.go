package types

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func makeInterval(s string) Interval {
	var i Interval
	if err := json.Unmarshal([]byte("\""+s+"\""), &i); err != nil {
		panic(err)
	}
	return i
}

func TestIntervalsCollide(t *testing.T) {
	for _, tc := range []struct {
		Name      string
		Intervals []Interval
		Expected  bool
	}{
		{
			Name:      "Empty interval slice",
			Intervals: nil,
			Expected:  false,
		},
		{
			Name: "Duplicate interval",
			Intervals: []Interval{
				makeInterval("12:00-14:00"),
				makeInterval("12:00-14:00"),
			},
			Expected: true,
		},
		{
			Name: "No collision",
			Intervals: []Interval{
				makeInterval("12:00-14:00"),
				makeInterval("15:00-17:00"),
				makeInterval("07:00-11:00"),
			},
			Expected: false,
		},
		{
			Name: "Collide by strict match",
			Intervals: []Interval{
				makeInterval("12:00-14:00"),
				makeInterval("14:00-16:00"),
			},
			Expected: true,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			require.Equal(t, tc.Expected, IntervalsCollide(tc.Intervals...))
		})
	}
}
