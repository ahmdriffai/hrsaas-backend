package lib

import (
	"strings"
	"time"
)

type DateOnly struct {
	time.Time
}

func (d *DateOnly) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" || s == "null" {
		return nil
	}

	t, err := time.Parse("2006-01-02", string(s))
	if err != nil {
		return err
	}

	d.Time = t
	return nil
}
