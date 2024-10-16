package auto

import (
	"encoding/json"
	"time"
)

const (
	timeFormatLayout = "2006-01-02 15:04:05"
)

type TimeLabel time.Time

func (t TimeLabel) MarshalJSON() ([]byte, error) {

	return json.Marshal((time.Time(t)).Format(timeFormatLayout))
}
