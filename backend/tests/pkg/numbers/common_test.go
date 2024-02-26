package numbers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandomInt(t *testing.T) {
	type args struct {
		min int
		max int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				min: 1,
				max: 1,
			},
			want: 1,
		},
		{
			name: "-1",
			args: args{
				min: -1,
				max: -1,
			},
			want: -1,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := RandomInt(tt.args.min, tt.args.max)
			assert.Equal(t, tt.want, got)
		})
	}
}
