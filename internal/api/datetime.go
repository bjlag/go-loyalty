package api

import "time"

type Datetime time.Time

func (d Datetime) MarshalJSON() ([]byte, error) {
	t := time.Time(d)
	formatted := t.Format(time.RFC3339)
	return []byte("\"" + formatted + "\""), nil
}
