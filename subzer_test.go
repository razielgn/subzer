package main

import (
    "testing"
    "time"
)

func TestParseTimestamp1(t *testing.T) {
    timestamp := "00:00:00,000"
    expected_duration := 0 * time.Second
    duration := ParseTimestamp(timestamp)

    if duration != expected_duration {
        t.Errorf("Expected duration to equal %s, got %s", expected_duration, duration)
    }
}

func TestParseTimestamp2(t *testing.T) {
    timestamp := "08:17:20,301"
    expected_duration := 8*time.Hour + 17*time.Minute + 20*time.Second + 301*time.Millisecond
    duration := ParseTimestamp(timestamp)

    if duration != expected_duration {
        t.Errorf("Expected duration to equal %s, got %s", expected_duration, duration)
    }
}

func TestReduceDurations1(t *testing.T) {
    start := ParseTimestamp("00:00:02,110")
    end := ParseTimestamp("00:00:04,450")

    seconds := ReduceDuration(start, end)
    if seconds != 2 {
        t.Errorf("Expected seconds to equal 2, got %d", seconds)
    }
}

func TestReduceDurations2(t *testing.T) {
    start := ParseTimestamp("00:00:00,110")
    end := ParseTimestamp("00:00:08,950")

    seconds := ReduceDuration(start, end)
    if seconds != 8 {
        t.Errorf("Expected seconds to equal 8, got %d", seconds)
    }
}
