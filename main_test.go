package main

import (
	"os"
	"reflect"
	"testing"
)

func TestLamportClock(t *testing.T) {
	// Create a temporary file with test data
	tmpfile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	//

	testData := "p0: s1 r1 l r1\np1: s0 s2 r0 l s2 s0 l r2\np2: l s1 r1 r1\n"
	if _, err := tmpfile.Write([]byte(testData)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	s := NewSimulator(tmpfile.Name())
	s.Simulate()

	expectedEvents := map[int][]Event{
		0: {
			{sent, 1, 0, 1},
			{recv, 2, 1, 0},
			{t: local, timestamp: 3},
			{recv, 7, 1, 0},
		},
		1: {
			{sent, 1, 1, 0},
			{sent, 2, 1, 2},
			{recv, 3, 0, 1},
			{t: local, timestamp: 4},
			{sent, 5, 1, 2},
			{sent, 6, 1, 0},
			{t: local, timestamp: 7},
			{recv, 8, 2, 1},
		},
		2: {
			{t: local, timestamp: 1},
			{sent, 2, 2, 1},
			{recv, 3, 1, 2},
			{recv, 6, 1, 2},
		},
	}
	eq := reflect.DeepEqual(expectedEvents, s.events)
	if !eq {
		t.Fatalf(
			"expected events map:\n%v\ndoesn't equal actual events map:\n%v",
			expectedEvents,
			s.events,
		)
	}
}
