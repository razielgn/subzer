package main

import "testing"

func TestTruth(t *testing.T) {
    if !truth() {
        t.Error("Expected the truth, liar.")
    }
}
