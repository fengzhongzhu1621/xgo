package datetime

import "time"

func Duration(str string) (time.Duration, error) {
	dur, err := time.ParseDuration(str)
	if err != nil {
		return time.Duration(0), err
	}
	return dur, nil
}
