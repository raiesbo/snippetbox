package main

import (
	"testing"
	"time"

	"github.com/raiesbo/snippetbox/internal/assert"
)

func TestHumanDate(t *testing.T) {
	// Create a slice of anonymous structs containing the test case name,
	// input to out humanDate() funciton (the tmb field), and expected output
	// (the want field).
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2024, 3, 17, 10, 15, 0, 0, time.UTC),
			want: "17 Mar 2024 at 10:15",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2024, 3, 17, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "17 Mar 2024 at 10:15",
		},
	}

	// Loop over the test cases.
	for _, tt := range tests {
		// use the t.Run() function to run a sub-test for each test case. The
		// first parameter to this is the name of the test (which is used to
		// identify the sub-test in any log output) and the second paramter is
		// and anonymous function containing the actual test for each case.
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)
			if hd != tt.want {
				t.Errorf("got %q; want %q", hd, tt.want)
			}
			// Use the new assert.Equal() helper to compare the expected and
			// actual values.
			assert.Equal(t, hd, tt.want)
		})
	}

	// // Initialize a new time.Time object and pass it to the humanDate function.
	// tm := time.Date(2024, 3, 17, 10, 15, 0, 0, time.UTC)
	// hd := humanDate(tm)

	// // Check that the output from the humanDate function is in the format we
	// // expect. If it isn't what we expect, use the t.Errorf() funciton to
	// // indicate that the test has failed and log the expected and actual values.
	// expectedResult := "17 Mar 2024 at 10:15"
	// if hd != expectedResult {
	// 	t.Errorf("got %q; want %q", hd, expectedResult)
	// }
}
