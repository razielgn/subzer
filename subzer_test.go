package main

import (
    stringio "github.com/kdar/stringio"
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
    start := ParseTimestamp("00:00:07,390")
    end := ParseTimestamp("00:00:09,280")

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

func TestSrtBlockEquality(t *testing.T) {
    start := (0 * time.Second)
    end := (2*time.Second + 110*time.Millisecond)
    text := []string{"[Salim Ismail] [Global Ambassador, Singularity University] Do you think you have what it takes"}
    expected := NewSrtBlock(1, start, end, text)
    output := NewSrtBlock(1, start, end, text)

    if expected.String() != output.String() {
        t.Errorf("Expected parsed srt block to be\n%#+v\ngot\n%#+v", expected, output)
    }
}

func SrtBlockPieces() []string {
    return []string{"1", "00:00:00,000 --> 00:00:02,110", "[Salim Ismail] [Global Ambassador, Singularity University] Do you think you have what it takes", "Whatever"}
}

func TestSrtBlockParsing1(t *testing.T) {
    srt_block := SrtBlockPieces()

    start := (0 * time.Second)
    end := (2*time.Second + 110*time.Millisecond)
    text := []string{"[Salim Ismail] [Global Ambassador, Singularity University] Do you think you have what it takes", "Whatever"}
    expected := NewSrtBlock(1, start, end, text)
    output := SrtBlockParse(srt_block)

    if expected.String() != output.String() {
        t.Errorf("Expected parsed srt block to be\n%#+v\ngot\n%#+v", expected, output)
    }
}

func TestSrtBlockParsing2(t *testing.T) {
    srt_block := SrtBlockPieces()

    start := (0 * time.Second)
    end := (2*time.Second + 110*time.Millisecond)
    text := []string{"[Sxlim Ismail] [Global Ambassador, Singularity University] Do you think you have what it takes"}
    expected := NewSrtBlock(1, start, end, text)
    output := SrtBlockParse(srt_block)

    if expected.String() == output.String() {
        t.Errorf("Expected parsed srt block to be\n%#v\ngot\n%#v", expected, output)
    }
}

func SourceStringExample() string {
    return `1
00:00:00,000 --> 00:00:02,110
[Salim Ismail] [Global Ambassador, Singularity University] Do you think you have what it takes

2
00:00:02,110 --> 00:00:04,450
to positively impact a billion people around the world?

3
00:00:04,450 --> 00:00:07,390
Here are a few of our students that do.

4
00:00:07,390 --> 00:00:09,280
[♪ Music ♪]

`
}

func DestinationStringExample() string {
    return "000009\twhatever.txt\n" +
        "000000\t[Salim Ismail] [Global Ambassador, Singularity University] Do you think you have what it takes\n" +
        "000002\tto positively impact a billion people around the world?\n" +
        "000004\tHere are a few of our students that do.\n" +
        "000007\t[♪ Music ♪]\n"
}

func TestConversion1(t *testing.T) {
    source := stringio.New()
    destination := stringio.New()

    source.WriteString(SourceStringExample())
    source.Seek(0, 0)

    StreamConversion("whatever.txt", source, destination, destination)
    output := destination.GetValueString()

    expected_output := DestinationStringExample()

    if expected_output != output {
        t.Errorf("Expected output to be\n%s\ngot\n\n%s\n", expected_output, output)
    }
}
