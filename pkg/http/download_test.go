package http

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsSuccess(t *testing.T) {
	type args struct {
		statusCode int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: http.StatusText(http.StatusOK),
			args: args{http.StatusOK},
			want: true,
		},
		{
			name: http.StatusText(http.StatusCreated),
			args: args{http.StatusCreated},
			want: true,
		},
		{
			name: http.StatusText(http.StatusIMUsed),
			args: args{http.StatusIMUsed},
			want: true,
		},
		{
			name: http.StatusText(http.StatusContinue),
			args: args{http.StatusContinue},
			want: false,
		},
		{
			name: http.StatusText(http.StatusMultipleChoices),
			args: args{http.StatusMultipleChoices},
			want: false,
		},
		{
			name: http.StatusText(http.StatusMovedPermanently),
			args: args{http.StatusMovedPermanently},
			want: false,
		},
		{
			name: http.StatusText(http.StatusBadRequest),
			args: args{http.StatusBadRequest},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsSuccess(tt.args.statusCode))
		})
	}
}
