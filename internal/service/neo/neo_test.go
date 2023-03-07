package neo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseDates(t *testing.T) {
	tests := []struct {
		name    string
		dates   []string
		want    []time.Time
		wantErr bool
	}{
		{
			name:    "Valid dates",
			dates:   []string{"2016-12-01"},
			want:    []time.Time{time.Date(2016, 12, 1, 0, 0, 0, 0, time.UTC)},
			wantErr: false,
		},
		{
			name:    "Invalid dates",
			dates:   []string{"invalid date"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Empty dates",
			dates:   []string{},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			got, err := ParseDates(tt.dates)

			assert.Equal(tt.wantErr, err != nil)

			assert.Equal(tt.want, got)
		})
	}
}
