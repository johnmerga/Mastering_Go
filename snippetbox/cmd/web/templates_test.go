package main

import (
	"testing"
	"time"

	"github.com/johnmerga/Mastering_Go/snippetbox/internal/assert"
)

func Test_HumanDate(t *testing.T) {
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2023, 12, 3, 9, 9, 49, 49, time.UTC),
			want: "03 Dec 2023 - 09:09",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2023, 12, 3, 9, 9, 49, 49, time.FixedZone("CET", 1)),
			want: "03 Dec 2023 - 09:09",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)
			// we can also use exclude [string]
			// assert.Equals(t, hd, tt.want)
			assert.Equals[string](t, hd, tt.want)
		})
	}
}
