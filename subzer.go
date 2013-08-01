package main

import (
    "regexp"
    "strconv"
    "time"
)

func ReduceDuration(start time.Duration, end time.Duration) int {
    return int(end.Seconds() - start.Seconds())
}

func ParseTimestamp(timestamp string) time.Duration {
    pieces := regexp.MustCompile(":|,").Split(timestamp, 4)

    hours, _ := strconv.ParseUint(pieces[0], 10, 8)
    minutes, _ := strconv.ParseUint(pieces[1], 10, 8)
    seconds, _ := strconv.ParseUint(pieces[2], 10, 8)
    milliseconds, _ := strconv.ParseUint(pieces[3], 10, 16)

    return time.Duration(hours)*time.Hour +
        time.Duration(minutes)*time.Minute +
        time.Duration(seconds)*time.Second +
        time.Duration(milliseconds)*time.Millisecond
}

func main() {
}
