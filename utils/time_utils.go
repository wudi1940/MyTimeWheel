package utils

import "time"

func TimeToSec(t time.Time) uint64 {
	return uint64(t.UnixNano() / int64(time.Second))
}

func SecToTime(t int64) time.Time {
	return time.Unix(0, t*int64(time.Second)).UTC()
}

// TimeToMs returns an integer number, which represents t in milliseconds.
func TimeToMs(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

// MsToTime returns the UTC time corresponding to the given Unix time,
// t milliseconds since January 1, 1970 UTC.
func MsToTime(t int64) time.Time {
	return time.Unix(0, t*int64(time.Millisecond)).UTC()
}
