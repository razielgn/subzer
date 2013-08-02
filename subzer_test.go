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
0:0:0,100 --> 0:0:8,000
Did you know?

2
0:0:9,100 --> 0:0:11,000
Did you know?

3
0:0:11,100 --> 0:0:14,000
in the next 8 seconds

4
0:0:15,100 --> 0:0:17,000
34 babies will be born

`
}

func DestinationStringExample() string {
    return "000017\twhatever.txt\n" +
        "000000\tDid you know?\n" +
        "000009\tDid you know?\n" +
        "000011\tin the next 8 seconds\n" +
        "000015\t34 babies will be born\n"
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
